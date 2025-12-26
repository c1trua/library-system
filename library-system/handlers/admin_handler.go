package handlers

import (
	"errors"
	"library-system/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminService *services.AdminService
}

func NewAdminHandler(adminService *services.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// AddBook godoc
// @Summary 添加图书
// @Description 管理员添加新图书到系统
// @Tags admin
// @Accept json
// @Produce json
// @Param request body AddBookRequest true "图书信息"
// @Success 200 {object} SuccessResponse "添加成功"
// @Failure 400 {object} ErrorResponse "请求参数错误或格式不正确"
// @Failure 409 {object} ErrorResponse "图书已存在"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /admin/books [post]
func (h *AdminHandler) AddBook(c *gin.Context) {
	var req AddBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数格式错误", err)
		return
	}

	err := h.adminService.AddBook(req.Title, req.Author, req.Stock)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			BadRequest(c, "请求参数错误", err)
			return
		} else if errors.Is(err, services.ErrBookExists) {
			Conflict(c, "图书已存在", err)
			return
		} else {
			InternalError(c, "添加失败", err)
			return
		}
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "图书添加成功"})
}

// UpdateBook godoc
// @Summary 更新图书信息
// @Description 管理员更新现有图书信息
// @Tags admin
// @Accept json
// @Produce json
// @Param request body UpdateBookRequest true "图书更新信息"
// @Success 200 {object} SuccessResponse "更新成功"
// @Failure 400 {object} ErrorResponse "请求参数错误或格式不正确"
// @Failure 404 {object} ErrorResponse "未找到该图书"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /admin/books [put]
func (h *AdminHandler) UpdateBook(c *gin.Context) {
	var req UpdateBookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数格式错误", err)
		return
	}

	err := h.adminService.UpdateBook(req.Title, req.Author, req.ID, req.Stock)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			BadRequest(c, "请求参数错误", err)
			return
		} else if errors.Is(err, services.ErrBookNotFound) {
			NotFound(c, "未找到该图书", err)
			return
		} else {
			InternalError(c, "更新失败", err)
			return
		}
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "图书更新成功"})
}

// DeleteBook godoc
// @Summary 删除图书
// @Description 管理员从系统中删除图书
// @Tags admin
// @Accept json
// @Produce json
// @Param request body DeleteBookRequest true "删除图书请求"
// @Success 200 {object} SuccessResponse "删除成功"
// @Failure 400 {object} ErrorResponse "请求参数错误或格式不正确"
// @Failure 404 {object} ErrorResponse "未找到该图书"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /admin/books [delete]
func (h *AdminHandler) DeleteBook(c *gin.Context) {
	var req DeleteBookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数格式错误", err)
		return
	}

	err := h.adminService.DeleteBook(req.ID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			BadRequest(c, "请求参数错误", err)
			return
		} else if errors.Is(err, services.ErrBookNotFound) {
			NotFound(c, "未找到该图书", err)
			return
		} else {
			InternalError(c, "删除失败", err)
			return
		}
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "图书删除成功"})
}

// GetAllBorrowRecords godoc
// @Summary 获取所有借阅记录
// @Description 管理员查看所有用户的借阅记录
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {array} models.BorrowRecord "借阅记录数组"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /admin/borrow-records [get]
func (h *AdminHandler) GetAllBorrowRecords(c *gin.Context) {
	records, err := h.adminService.GetAllBorrowRecords()
	if err != nil {
		InternalError(c, "获取借阅记录失败", err)
		return
	}

	c.JSON(http.StatusOK, records)
}

// 请求和响应结构体定义
type AddBookRequest struct {
	Title  string `json:"title" binding:"required" example:"LemonisTheBestFruit"`
	Author string `json:"author" binding:"required" example:"Lemon"`
	Stock  int    `json:"stock" binding:"required" example:"10"`
}

type UpdateBookRequest struct {
	ID     int    `json:"id" binding:"required" example:"1"`
	Title  string `json:"title" binding:"required" example:"LemonisTheBestFruit"`
	Author string `json:"author" binding:"required" example:"Lemon"`
	Stock  int    `json:"stock" binding:"required" example:"15"`
}

type DeleteBookRequest struct {
	ID int `json:"id" binding:"required" example:"1"`
}

type SuccessResponse struct {
	Message string `json:"message" example:"操作成功"`
}
