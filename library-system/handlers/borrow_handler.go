package handlers

import (
	"errors"
	"library-system/models"
	"library-system/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BorrowHandler struct {
	borrowService *services.BorrowService
}

func NewBorrowHandler(borrowService *services.BorrowService) *BorrowHandler {
	return &BorrowHandler{borrowService: borrowService}
}

// BorrowBook godoc
// @Summary 借阅图书
// @Description 用户借阅指定图书（需要登录）
// @Tags borrow
// @Accept json
// @Produce json
// @Param request body BorrowBookRequest true "借书信息"
// @Success 200 {object} SuccessResponse "借书成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 401 {object} ErrorResponse "用户未认证"
// @Failure 404 {object} ErrorResponse "图书不存在"
// @Failure 409 {object} ErrorResponse "库存不足或借阅次数已达上限"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /borrow [post]
func (h *BorrowHandler) BorrowBook(c *gin.Context) {
	var req BorrowBookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数格式错误", err)
		return
	}

	// 获取用户信息
	userObj, exists := c.Get("user")
	if !exists {
		Unauthorized(c, "未找到用户信息", nil)
		return
	}

	user := userObj.(*models.User)

	// 借书
	err := h.borrowService.BorrowBook(user.ID, req.BookID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			BadRequest(c, "请求参数错误", err)
			return
		} else if errors.Is(err, services.ErrBookNotFound) {
			NotFound(c, "未找到该图书", err)
			return
		} else if errors.Is(err, services.ErrStockNotEnough) {
			Conflict(c, "库存不足", err)
			return
		} else if errors.Is(err, services.ErrBorrowLimit) {
			Conflict(c, "借阅次数已达上限", err)
			return
		} else {
			InternalError(c, "借阅失败", err)
			return
		}
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "借书成功"})
}

// ReturnBook godoc
// @Summary 归还图书
// @Description 用户归还已借阅的图书（需要登录）
// @Tags borrow
// @Accept json
// @Produce json
// @Param request body ReturnBookRequest true "还书信息"
// @Success 200 {object} SuccessResponse "还书成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 401 {object} ErrorResponse "用户未认证"
// @Failure 403 {object} ErrorResponse "权限不足"
// @Failure 404 {object} ErrorResponse "借阅记录或图书不存在"
// @Failure 409 {object} ErrorResponse "图书已归还"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /borrow/return [post]
func (h *BorrowHandler) ReturnBook(c *gin.Context) {
	var req ReturnBookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数格式错误", err)
		return
	}

	// 获取用户信息
	userObj, exists := c.Get("user")
	if !exists {
		Unauthorized(c, "未找到用户信息", nil)
		return
	}
	user := userObj.(*models.User)

	// 还书
	err := h.borrowService.ReturnBook(req.RecordID, user.ID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			BadRequest(c, "请求参数错误", err)
			return
		} else if errors.Is(err, services.ErrRecordNotFound) {
			NotFound(c, "未找到该借阅记录", err)
			return
		} else if errors.Is(err, services.ErrPermissionDenied) {
			Forbidden(c, "借阅者与当前用户不匹配", err)
			return
		} else if errors.Is(err, services.ErrAlreadyReturned) {
			Conflict(c, "图书已归还", err)
			return
		} else if errors.Is(err, services.ErrBookNotFound) {
			NotFound(c, "未找到该图书", err)
			return
		} else {
			InternalError(c, "还书失败", err)
			return
		}
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "还书成功"})
}

// GetUserBorrowRecords godoc
// @Summary 获取用户借阅记录
// @Description 获取当前用户的所有借阅记录（需要登录）
// @Tags borrow
// @Accept json
// @Produce json
// @Success 200 {array} models.BorrowRecord "借阅记录数组"
// @Failure 401 {object} ErrorResponse "用户未认证"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /borrow/records [get]
func (h *BorrowHandler) GetUserBorrowRecords(c *gin.Context) {
	// 获取用户信息
	userObj, exists := c.Get("user")
	if !exists {
		Unauthorized(c, "未找到用户信息", nil)
		return
	}
	user := userObj.(*models.User)

	// 获取借阅记录
	records, err := h.borrowService.GetUserBorrowRecords(user.ID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			BadRequest(c, "请求参数错误", err)
			return
		} else if errors.Is(err, services.ErrRecordNotFound) {
			c.JSON(http.StatusOK, []models.BorrowRecord{})
			return
		} else {
			InternalError(c, "获取借阅记录失败", err)
			return
		}
	}

	c.JSON(http.StatusOK, records)
}

// 请求和响应结构体定义
type BorrowBookRequest struct {
	BookID int `json:"book_id" binding:"required" example:"1"`
}

type ReturnBookRequest struct {
	RecordID int `json:"record_id" binding:"required" example:"1"`
}
