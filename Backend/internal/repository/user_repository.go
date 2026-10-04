package repository

import (
	"Backend/internal/model"
	"errors"

	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	FindbyUsername(username string) (*model.User, error)
	FindById(id uint) (*model.User, error)
	List(limit, offset int) ([]model.User, int64, error)
	Create(u *model.User) error
	UpdateRoleSatker(id int, role string, satkerID *int) error
	SetActive(id int, active bool) error
	UpdatePassword(id int, passwordHash string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindbyUsername(username string) (*model.User, error) {
	var user model.User

	err := r.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

func (r *userRepository) FindById(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

func (r *userRepository) List(limit, offset int) ([]model.User, int64, error) {
	var total int64
	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var users []model.User
	err := r.db.Order("id ASC").Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

func (r *userRepository) Create(u *model.User) error {
	return r.db.Create(u).Error
}

// checkAffected: update ke id yang tidak ada → ErrUserNotFound.
func checkAffected(res *gorm.DB) error {
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userRepository) UpdateRoleSatker(id int, role string, satkerID *int) error {
	return checkAffected(r.db.Model(&model.User{}).Where("id = ?", id).
		Updates(map[string]any{"role": role, "satker_id": satkerID}))
}

func (r *userRepository) SetActive(id int, active bool) error {
	return checkAffected(r.db.Model(&model.User{}).Where("id = ?", id).
		Update("is_active", active))
}

func (r *userRepository) UpdatePassword(id int, passwordHash string) error {
	return checkAffected(r.db.Model(&model.User{}).Where("id = ?", id).
		Update("password_hash", passwordHash))
}
