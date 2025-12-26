package handlers

import (
	"errors"
	"library-system/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	bookService *services.BookService
}

func NewBookHandler(bookService *services.BookService) *BookHandler {
	return &BookHandler{bookService: bookService}
}

// GetAllBooks godoc
// @Summary 获取所有图书
// @Description 获取系统中的所有图书列表
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {array} models.Book "图书列表" // 修改1：{array} 改为 {object}，因为返回的是单个模型实例的列表
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /books [get]
func (h *BookHandler) GetAllBooks(c *gin.Context) {
	books, err := h.bookService.GetAllBooks()
	if err != nil {
		InternalError(c, "无法获取图书列表", err)
		return
	}

	c.JSON(http.StatusOK, books)
}

// GetBookInfoByID godoc
// @Summary 根据ID获取图书信息
// @Description 通过图书ID获取特定图书的详细信息
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "图书ID"
// @Success 200 {object} models.Book "图书信息" // 修改1：{array} 改为 {object}，返回单个图书对象
// @Failure 400 {object} ErrorResponse "无效的图书ID"
// @Failure 404 {object} ErrorResponse "图书不存在"
// @Failure 500 {object} ErrorResponse "服务器内部错误" // 修改3：补充500错误
// @Router /books/{id} [get]
func (h *BookHandler) GetBookInfoByID(c *gin.Context) {
	// 从路径参数获取ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		BadRequest(c, "无效的图书ID", err)
		return
	}

	// 获取图书信息
	book, err := h.bookService.GetBookInfoByID(id)
	if err != nil {
		if errors.Is(err, services.ErrBookNotFound) {
			NotFound(c, "未找到该图书", err)
			return
		} else {
			InternalError(c, "获取图书信息失败", err)
			return
		}
	}

	c.JSON(http.StatusOK, book)
}

// GetBookInfoByTitle godoc
// @Summary 根据标题获取图书信息
// @Description 通过图书标题获取特定图书的详细信息
// @Tags books
// @Accept json
// @Produce json
// @Param title query string true "图书标题"
// @Success 200 {object} models.Book "图书信息"
// @Failure 400 {object} ErrorResponse "标题不能为空"
// @Failure 404 {object} ErrorResponse "图书不存在"
// @Failure 500 {object} ErrorResponse "服务器内部错误" // 修改2：补充500错误
// @Router /books/title [get]
func (h *BookHandler) GetBookInfoByTitle(c *gin.Context) {
	// 从查询参数获取标题
	title := c.Query("title")
	if title == "" {
		BadRequest(c, "图书标题不能为空", nil)
		return
	}

	// 获取图书信息
	book, err := h.bookService.GetBookInfoByTitle(title)
	if err != nil {
		if errors.Is(err, services.ErrBookNotFound) {
			NotFound(c, "未找到该图书", err)
			return
		} else {
			InternalError(c, "获取图书信息失败", err)
			return
		}
	}

	c.JSON(http.StatusOK, book)
}

// SearchBooksByKeyword godoc
// @Summary 根据关键词搜索图书
// @Description 通过关键词搜索图书，支持标题、作者等字段的模糊匹配
// @Tags books
// @Accept json
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Success 200 {array} models.Book "搜索结果" // 此处使用 {array} 正确，因为返回的是图书列表
// @Failure 400 {object} ErrorResponse "搜索关键词不能为空"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /books/search [get]
func (h *BookHandler) SearchBooksByKeyword(c *gin.Context) {
	// 从查询参数获取关键词
	keyword := c.Query("keyword")
	if keyword == "" {
		BadRequest(c, "搜索关键词不能为空", nil)
		return
	}

	// 搜索图书
	books, err := h.bookService.SearchBooksByKeyword(keyword)
	if err != nil {
		InternalError(c, "无法搜索图书", err)
		return
	}

	c.JSON(http.StatusOK, books)
}

// SearchBooksByTitleKeyword godoc
// @Summary 根据书名关键词搜索图书
// @Description 通过书名中的关键词进行模糊搜索，返回匹配的图书列表
// @Tags books
// @Accept json
// @Produce json
// @Param titlekeyword query string true "书名关键词"
// @Success 200 {array} models.Book "搜索到的图书列表"
// @Failure 400 {object} ErrorResponse "书名关键词不能为空"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /books/search/title [get]
func (h *BookHandler) SearchBooksByTitleKeyword(c *gin.Context) {
	// 从查询参数获取标题
	titlekeyword := c.Query("titlekeyword")
	if titlekeyword == "" {
		BadRequest(c, "图书标题关键词不能为空", nil)
		return
	}

	// 搜索图书
	books, err := h.bookService.SearchBooksByTitleKeyword(titlekeyword)
	if err != nil {
		InternalError(c, "无法搜索图书", err)
		return
	}

	c.JSON(http.StatusOK, books)
}

// SearchBooksByAuthor godoc
// @Summary 根据作者搜索图书
// @Description 通过精确的作者名称搜索图书，返回该作者的所有图书
// @Tags books
// @Accept json
// @Produce json
// @Param author query string true "作者名称"
// @Success 200 {array} models.Book "搜索到的图书列表"
// @Failure 400 {object} ErrorResponse "作者不能为空"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /books/search/author [get]
func (h *BookHandler) SearchBooksByAuthor(c *gin.Context) {
	// 从查询参数获取作者
	author := c.Query("author")
	if author == "" {
		BadRequest(c, "作者不能为空", nil)
		return
	}

	// 搜索图书
	books, err := h.bookService.SearchBooksByAuthor(author)
	if err != nil {
		InternalError(c, "无法搜索图书", err)
		return
	}

	c.JSON(http.StatusOK, books)
}
