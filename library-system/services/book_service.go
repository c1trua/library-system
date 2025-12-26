package services

import (
	"errors"
	"fmt"
	"library-system/models"
	"library-system/repositories"

	"gorm.io/gorm"
)

type BookService struct {
	bookRepo repositories.BookRepository
}

func NewBookService(bookRepo repositories.BookRepository) *BookService {
	return &BookService{bookRepo: bookRepo}
}

// GetAllBooks
func (s *BookService) GetAllBooks() ([]*models.Book, error) {
	books, err := s.bookRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all books: %w", err)
	}

	return books, nil
}

// GetBookInfoByID
func (s *BookService) GetBookInfoByID(bookID int) (*models.Book, error) {
	book, err := s.bookRepo.GetByID(bookID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBookNotFound
		}
		return nil, fmt.Errorf("failed to get book by ID: %w", err)
	}

	return book, nil
}

// GetBookInfoByTitle
func (s *BookService) GetBookInfoByTitle(title string) (*models.Book, error) {
	book, err := s.bookRepo.GetByTitle(title)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBookNotFound
		}
		return nil, fmt.Errorf("failed to get book by title: %w", err)
	}

	return book, nil
}

// SearchBooksByKeyword
func (s *BookService) SearchBooksByKeyword(keyword string) ([]*models.Book, error) {
	books, err := s.bookRepo.SearchByKeyword(keyword)
	if err != nil {
		return []*models.Book{}, fmt.Errorf("failed to search books by keyword: %w", err)
	}

	return books, nil
}

// SearchBooksByTitleKeyword
func (s *BookService) SearchBooksByTitleKeyword(titlekeyword string) ([]*models.Book, error) {
	books, err := s.bookRepo.SearchByTitleKeyword(titlekeyword)
	if err != nil {
		return []*models.Book{}, fmt.Errorf("failed to search books by title keyword: %w", err)
	}

	return books, nil
}

// SearchBooksByAuthor
func (s *BookService) SearchBooksByAuthor(author string) ([]*models.Book, error) {
	books, err := s.bookRepo.SearchByAuthor(author)
	if err != nil {
		return []*models.Book{}, fmt.Errorf("failed to search books by author: %w", err)
	}

	return books, nil
}
