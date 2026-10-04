package repository

import (
	"Backend/internal/model"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrRiwayatNotFound        = errors.New("riwayat jabatan tidak ditemukan")
	ErrTmtSebelumJabatanAktif = errors.New("tanggal mulai jabatan baru harus setelah tanggal mulai jabatan aktif saat ini")
)

type RiwayatJabatanRepository interface {
	ListByPersonel(personelID int64) ([]model.RiwayatJabatan, error)
	FindByID(id int64) (*model.RiwayatJabatan, error)
	Create(rj *model.RiwayatJabatan) error
	Update(rj *model.RiwayatJabatan) error
	SoftDelete(id int64) error
}

type riwayatJabatanRepository struct {
	db *gorm.DB
}

func NewRiwayatJabatanRepository(db *gorm.DB) RiwayatJabatanRepository {
	return &riwayatJabatanRepository{db: db}
}

func withRelasi(db *gorm.DB) *gorm.DB {
	return db.Preload("Satker").Preload("Fungsi").Preload("Nivelering")
}

// ListByPersonel: seluruh riwayat, kronologis (lama → baru).
func (r *riwayatJabatanRepository) ListByPersonel(personelID int64) ([]model.RiwayatJabatan, error) {
	var items []model.RiwayatJabatan
	err := withRelasi(r.db).
		Where("personel_id = ?", personelID).
		Order("tmt_mulai ASC").
		Find(&items).Error
	return items, err
}

func (r *riwayatJabatanRepository) FindByID(id int64) (*model.RiwayatJabatan, error) {
	var rj model.RiwayatJabatan
	err := withRelasi(r.db).First(&rj, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRiwayatNotFound
	}
	if err != nil {
		return nil, err
	}
	return &rj, nil
}

// Create: kalau jabatan baru DEFINITIF & masih aktif (tmt_selesai kosong),
// jabatan DEFINITIF aktif sebelumnya otomatis ditutup per H-1 tmt_mulai baru.
// Semua dalam 1 transaksi — gagal satu, batal semua.
func (r *riwayatJabatanRepository) Create(rj *model.RiwayatJabatan) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if rj.StatusJabatan == "DEFINITIF" && rj.TmtSelesai == nil {
			var aktif model.RiwayatJabatan
			err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("personel_id = ? AND tmt_selesai IS NULL AND status_jabatan = ?", rj.PersonelID, "DEFINITIF").
				First(&aktif).Error

			switch {
			case err == nil:
				if !rj.TmtMulai.After(aktif.TmtMulai.Time) {
					return ErrTmtSebelumJabatanAktif
				}
				selesai := model.Date{Time: rj.TmtMulai.AddDate(0, 0, -1)}
				if err := tx.Model(&model.RiwayatJabatan{}).
					Where("id = ?", aktif.ID).
					Update("tmt_selesai", selesai).Error; err != nil {
					return err
				}
			case !errors.Is(err, gorm.ErrRecordNotFound):
				return err
			}
		}
		return tx.Omit(clause.Associations).Create(rj).Error
	})
}

// Update: kolom eksplisit, relasi tidak ikut disimpan.
func (r *riwayatJabatanRepository) Update(rj *model.RiwayatJabatan) error {
	return r.db.Model(&model.RiwayatJabatan{ID: rj.ID}).
		Select("nama_jabatan", "satker_id", "fungsi_id", "nivelering_id",
			"tmt_mulai", "tmt_selesai", "status_jabatan", "nomor_skep", "keterangan").
		Omit(clause.Associations).
		Updates(rj).Error
}

func (r *riwayatJabatanRepository) SoftDelete(id int64) error {
	return r.db.Delete(&model.RiwayatJabatan{}, id).Error
}
