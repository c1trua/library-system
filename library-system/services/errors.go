package services

import "errors"

var (
	ErrUserNotFound     = errors.New("用户不存在")
	ErrUserExists       = errors.New("用户已存在")
	ErrInvalidPassword  = errors.New("密码错误")
	ErrBookNotFound     = errors.New("图书不存在")
	ErrBookExists       = errors.New("图书已存在")
	ErrStockNotEnough   = errors.New("库存不足")
	ErrBorrowLimit      = errors.New("借书数量已达上限")
	ErrRecordNotFound   = errors.New("借阅记录不存在")
	ErrAlreadyReturned  = errors.New("图书已归还")
	ErrPermissionDenied = errors.New("权限不足")
	ErrInvalidInput     = errors.New("无效的输入参数")
)
