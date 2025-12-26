package repositories

import (
	"library-system/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByUserID(id int) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

// Create
func (r *userRepositoryImpl) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByUserID
func (r *userRepositoryImpl) GetByUserID(id int) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	return &user, result.Error
}

// GetByUsername
func (r *userRepositoryImpl) GetByUsername(username string) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "name = ?", username)
	return &user, result.Error
}
