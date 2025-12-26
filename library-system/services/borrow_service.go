package services

import (
	"errors"
	"fmt"
	"library-system/models"
	"library-system/repositories"
	"time"

	"gorm.io/gorm"
)

type BorrowService struct {
	db *gorm.DB
}

func NewBorrowService(db *gorm.DB) *BorrowService {
	return &BorrowService{
		db: db,
	}
}

// BorrowBook
func (s *BorrowService) BorrowBook(userID int, bookID int) error {
	// 参数基础校验
	if userID <= 0 || bookID <= 0 {
		return ErrInvalidInput
	}

	// 事务处理
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 创建仓库实例
		txBookRepo := repositories.NewBookRepository(tx)
		txRecordRepo := repositories.NewBorrowRecordRepository(tx)

		// 查找图书
		book, err := txBookRepo.GetByID(bookID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrBookNotFound
			}
			return fmt.Errorf("failed to get book by ID: %w", err)
		}

		// 检查图书是否可借
		if book.Stock <= 0 {
			return ErrStockNotEnough
		}

		// 检查用户借书是否已达上限
		Count, err := txRecordRepo.CountActiveBorrowsByUserID(userID)
		if err != nil {
			return fmt.Errorf("failed to count active borrows by user ID: %w", err)
		}
		if Count >= 5 {
			return ErrBorrowLimit
		}

		// 库存-1
		book.Stock--
		if err := txBookRepo.Update(book); err != nil {
			return fmt.Errorf("failed to update book stock: %w", err)
		}

		// 创建新记录
		newRecord := &models.BorrowRecord{
			UserID:     userID,
			BookID:     bookID,
			BorrowedAt: time.Now(),
			DueDate:    time.Now().AddDate(0, 1, 0),
		}
		if err := txRecordRepo.Create(newRecord); err != nil {
			return fmt.Errorf("failed to create borrow record: %w", err)
		}

		return nil
	})
}

// ReturnBook
func (s *BorrowService) ReturnBook(recordID int, currentUserID int) error {
	// 参数基础校验
	if recordID <= 0 || currentUserID <= 0 {
		return ErrInvalidInput
	}

	// 事务处理
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 创建仓库实例
		txBookRepo := repositories.NewBookRepository(tx)
		txRecordRepo := repositories.NewBorrowRecordRepository(tx)

		// 查找记录
		record, err := txRecordRepo.GetByID(recordID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrRecordNotFound
			}
			return fmt.Errorf("failed to get borrow record by ID: %w", err)
		}

		// 检查权限
		if record.UserID != currentUserID {
			return ErrPermissionDenied
		}

		// 检查记录是否已经归还
		if record.ReturnedAt != nil {
			return ErrAlreadyReturned
		}

		// 更新图书库存
		book, err := txBookRepo.GetByID(record.BookID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrBookNotFound
			}
			return fmt.Errorf("failed to get book by ID: %w", err)
		}

		book.Stock++
		if err := txBookRepo.Update(book); err != nil {
			return fmt.Errorf("failed to update book stock: %w", err)
		}

		// 更新借阅记录
		currentTime := time.Now()
		record.ReturnedAt = &currentTime
		if err := txRecordRepo.Update(record); err != nil {
			return fmt.Errorf("failed to update borrow record: %w", err)
		}

		return nil
	})
}

// GetUserBorrowRecords
func (s *BorrowService) GetUserBorrowRecords(userID int) ([]*models.BorrowRecord, error) {
	// 参数基础校验
	if userID <= 0 {
		return nil, ErrInvalidInput
	}

	// 创建仓库实例
	RecordRepo := repositories.NewBorrowRecordRepository(s.db)

	records, err := RecordRepo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*models.BorrowRecord{}, ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get borrow records by user ID: %w", err)
	}

	return records, nil
}
