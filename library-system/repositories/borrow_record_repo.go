package repositories

import (
	"library-system/models"

	"gorm.io/gorm"
)

type BorrowRecordRepository interface {
	Create(record *models.BorrowRecord) error
	Update(record *models.BorrowRecord) error
	GetByID(id int) (*models.BorrowRecord, error)
	GetByUserID(userID int) ([]*models.BorrowRecord, error)
	GetByBookID(bookID int) ([]*models.BorrowRecord, error)
	CountActiveBorrowsByUserID(userID int) (int64, error)
	GetAll() ([]*models.BorrowRecord, error)
}

type borrowRecordRepoImpl struct {
	db *gorm.DB
}

func NewBorrowRecordRepository(db *gorm.DB) BorrowRecordRepository {
	return &borrowRecordRepoImpl{db: db}
}

// Create
func (r *borrowRecordRepoImpl) Create(record *models.BorrowRecord) error {
	return r.db.Create(record).Error
}

// Update
func (r *borrowRecordRepoImpl) Update(record *models.BorrowRecord) error {
	return r.db.Save(record).Error
}

// GetByID
func (r *borrowRecordRepoImpl) GetByID(id int) (*models.BorrowRecord, error) {
	var record models.BorrowRecord
	result := r.db.First(&record, id)
	return &record, result.Error
}

// GetByUserID
func (r *borrowRecordRepoImpl) GetByUserID(userid int) ([]*models.BorrowRecord, error) {
	var records []*models.BorrowRecord
	result := r.db.Where("user_id = ?", userid).Find(&records)
	return records, result.Error
}

// GetByBookID
func (r *borrowRecordRepoImpl) GetByBookID(bookid int) ([]*models.BorrowRecord, error) {
	var records []*models.BorrowRecord
	result := r.db.Where("book_id = ?", bookid).Find(&records)
	return records, result.Error
}

// CountActiveBorrowsByUserID
func (r *borrowRecordRepoImpl) CountActiveBorrowsByUserID(userID int) (int64, error) {
	var count int64
	result := r.db.Model(&models.BorrowRecord{}).Where("user_id = ? AND returned_at IS NULL", userID).Count(&count)
	return count, result.Error
}

// GetAll
func (r *borrowRecordRepoImpl) GetAll() ([]*models.BorrowRecord, error) {
	var records []*models.BorrowRecord
	result := r.db.Find(&records)
	return records, result.Error
}
