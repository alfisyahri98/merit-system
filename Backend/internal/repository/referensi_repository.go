package repository

import (
	"Backend/internal/model"

	"gorm.io/gorm"
)

// ReferensiRepository: data master untuk pilihan di form.
type ReferensiRepository interface {
	Pangkat() ([]model.Pangkat, error)
	Satker(scopeSatkerID *int) ([]model.Satker, error)
	Fungsi() ([]model.Fungsi, error)
	Nivelering() ([]model.Nivelering, error)
}

type referensiRepository struct {
	db *gorm.DB
}

func NewReferensiRepository(db *gorm.DB) ReferensiRepository {
	return &referensiRepository{db: db}
}

func (r *referensiRepository) Pangkat() ([]model.Pangkat, error) {
	var items []model.Pangkat
	err := r.db.Order("jenis, urutan").Find(&items).Error
	return items, err
}

// Satker: operator hanya mendapat satker di wilayahnya.
func (r *referensiRepository) Satker(scopeSatkerID *int) ([]model.Satker, error) {
	q := r.db.Order("kode")
	if scopeSatkerID != nil {
		q = q.Where("id IN (SELECT id FROM satker_turunan(?))", *scopeSatkerID)
	}
	var items []model.Satker
	err := q.Find(&items).Error
	return items, err
}

func (r *referensiRepository) Fungsi() ([]model.Fungsi, error) {
	var items []model.Fungsi
	err := r.db.Order("nama").Find(&items).Error
	return items, err
}

func (r *referensiRepository) Nivelering() ([]model.Nivelering, error) {
	var items []model.Nivelering
	err := r.db.Order("urutan").Find(&items).Error
	return items, err
}
