package services

import (
	"errors"
	"fmt"
	"library-system/models"
	"library-system/repositories"
	"time"

	"gorm.io/gorm"
)

type AdminService struct {
	db *gorm.DB
}

func NewAdminService(db *gorm.DB) *AdminService {
	return &AdminService{
		db: db,
	}
}

// AddBook
func (s *AdminService) AddBook(title, author string, stock int) error {
	// 参数基础校验
	if title == "" || author == "" || stock < 0 {
		return ErrInvalidInput
	}

	// 创建仓库实例
	bookRepo := repositories.NewBookRepository(s.db)

	// 判断图书是否已存在
	_, err := bookRepo.GetByTitle(title)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check book existence: %w", err)
	}
	if err == nil {
		return ErrBookExists
	}

	book := &models.Book{
		Title:  title,
		Author: author,
		Stock:  stock,
	}

	if err := bookRepo.Create(book); err != nil {
		return fmt.Errorf("failed to create book: %w", err)
	}

	return nil
}

// UpdateBook
func (s *AdminService) UpdateBook(title, author string, ID, stock int) error {
	// 参数基础校验
	if title == "" || author == "" || ID < 0 || stock < 0 {
		return ErrInvalidInput
	}

	// 创建仓库实例
	bookRepo := repositories.NewBookRepository(s.db)

	// 查询图书
	book, err := bookRepo.GetByID(ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBookNotFound
		}
		return fmt.Errorf("failed to get book by ID: %w", err)
	}

	book.Title = title
	book.Author = author
	book.Stock = stock

	if err := bookRepo.Update(book); err != nil {
		return fmt.Errorf("failed to update book: %w", err)
	}

	return nil
}

// DeleteBook
func (s *AdminService) DeleteBook(ID int) error {
	// 参数基础校验
	if ID < 0 {
		return ErrInvalidInput
	}

	// 事务处理
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 创建仓库实例
		txBookRepo := repositories.NewBookRepository(tx)
		txRecordRepo := repositories.NewBorrowRecordRepository(tx)

		// 查询图书
		book, err := txBookRepo.GetByID(ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrBookNotFound
			}
			return fmt.Errorf("failed to get book by ID: %w", err)
		}

		// 将所有借书记录改为已归还
		borrowRecords, err := txRecordRepo.GetByBookID(ID)
		if err != nil {
			return fmt.Errorf("failed to get borrow records by book ID: %w", err)
		}

		for _, record := range borrowRecords {
			currentTime := time.Now()
			if record.ReturnedAt == nil {
				record.ReturnedAt = &currentTime

				if err := txRecordRepo.Update(record); err != nil {
					return fmt.Errorf("failed to update borrow record: %w", err)
				}
			}
		}

		if err := txBookRepo.Delete(book); err != nil {
			return fmt.Errorf("failed to delete book: %w", err)
		}

		return nil
	})
}

// GetAllBorrowRecords
func (s *AdminService) GetAllBorrowRecords() ([]*models.BorrowRecord, error) {
	// 创建仓库实例
	recordRepo := repositories.NewBorrowRecordRepository(s.db)

	records, err := recordRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all borrow records: %w", err)
	}

	return records, nil
}
