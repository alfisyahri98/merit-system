package repository

import (
	"Backend/internal/model"
	"errors"

	"gorm.io/gorm"
)

var ErrClientNotFound = errors.New("api client tidak ditemukan")

type ApiClientRepository interface {
	FindByClientID(clientID string) (*model.ApiClient, error)
	FindByID(id int) (*model.ApiClient, error)
	List(limit, offset int) ([]model.ApiClient, int64, error)
	Create(c *model.ApiClient) error
	Update(c *model.ApiClient) error
	SetActive(id int, active bool) error
	UpdateSecret(id int, secretHash string) error
}

type apiClientRepository struct {
	db *gorm.DB
}

func NewApiClientRepository(db *gorm.DB) ApiClientRepository {
	return &apiClientRepository{db: db}
}

func (r *apiClientRepository) FindByClientID(clientID string) (*model.ApiClient, error) {
	var c model.ApiClient
	err := r.db.Where("client_id = ?", clientID).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrClientNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *apiClientRepository) FindByID(id int) (*model.ApiClient, error) {
	var c model.ApiClient
	err := r.db.First(&c, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrClientNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *apiClientRepository) List(limit, offset int) ([]model.ApiClient, int64, error) {
	var total int64
	if err := r.db.Model(&model.ApiClient{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []model.ApiClient
	err := r.db.Order("id ASC").Limit(limit).Offset(offset).Find(&items).Error
	return items, total, err
}

func (r *apiClientRepository) Create(c *model.ApiClient) error {
	return r.db.Create(c).Error
}

// Update: hanya nama, scope, wilayah. client_id & secret tidak ikut berubah.
func (r *apiClientRepository) Update(c *model.ApiClient) error {
	res := r.db.Model(&model.ApiClient{ID: c.ID}).
		Select("nama_aplikasi", "scopes", "akses_satker_id").
		Updates(c)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrClientNotFound
	}
	return nil
}

func (r *apiClientRepository) SetActive(id int, active bool) error {
	res := r.db.Model(&model.ApiClient{}).Where("id = ?", id).Update("is_active", active)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrClientNotFound
	}
	return nil
}

func (r *apiClientRepository) UpdateSecret(id int, secretHash string) error {
	res := r.db.Model(&model.ApiClient{}).Where("id = ?", id).Update("secret_hash", secretHash)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrClientNotFound
	}
	return nil
}
