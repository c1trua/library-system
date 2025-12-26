package repositories

import (
	"library-system/models"

	"gorm.io/gorm"
)

type BookRepository interface {
	GetAll() ([]*models.Book, error)
	GetByID(id int) (*models.Book, error)
	GetByTitle(title string) (*models.Book, error)
	Create(book *models.Book) error
	Update(book *models.Book) error
	Delete(book *models.Book) error
	SearchByKeyword(keyword string) ([]*models.Book, error)
	SearchByTitleKeyword(title string) ([]*models.Book, error)
	SearchByAuthor(author string) ([]*models.Book, error)
}

type bookRepositoryImpl struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepositoryImpl{db: db}
}

// GetAll
func (r *bookRepositoryImpl) GetAll() ([]*models.Book, error) {
	var books []*models.Book
	result := r.db.Find(&books)
	return books, result.Error
}

// GetByID
func (r *bookRepositoryImpl) GetByID(id int) (*models.Book, error) {
	var book models.Book
	result := r.db.First(&book, id)
	return &book, result.Error
}

// GetByTitle
func (r *bookRepositoryImpl) GetByTitle(title string) (*models.Book, error) {
	var book models.Book
	result := r.db.First(&book, "title =?", title)
	return &book, result.Error
}

// Create
func (r *bookRepositoryImpl) Create(book *models.Book) error {
	return r.db.Create(book).Error
}

// Update
func (r *bookRepositoryImpl) Update(book *models.Book) error {
	return r.db.Save(book).Error
}

// Delete
func (r *bookRepositoryImpl) Delete(book *models.Book) error {
	return r.db.Delete(book).Error
}

// SearchByKeyword
func (r *bookRepositoryImpl) SearchByKeyword(keyword string) ([]*models.Book, error) {
	var books []*models.Book
	result := r.db.Where("title LIKE ? OR author LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Find(&books)
	return books, result.Error
}

// SearchByTitleKeyword
func (r *bookRepositoryImpl) SearchByTitleKeyword(titlekeyword string) ([]*models.Book, error) {
	var books []*models.Book
	result := r.db.Where("title LIKE ?", "%"+titlekeyword+"%").Find(&books)
	return books, result.Error
}

// SearchByAuthor
func (r *bookRepositoryImpl) SearchByAuthor(author string) ([]*models.Book, error) {
	var books []*models.Book
	result := r.db.Where("author = ?", author).Find(&books)
	return books, result.Error
}
