package handlers

import (
	"errors"
	"library-system/models"
	"library-system/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type AuthHandler struct {
	authService  *services.AuthService
	sessionStore sessions.Store
}

func NewAuthHandler(authService *services.AuthService, sessionStore sessions.Store) *AuthHandler {
	return &AuthHandler{authService: authService, sessionStore: sessionStore}
}

// Login godoc
// @Summary 用户登录
// @Description 用户使用用户名和密码登录系统，登录成功后设置Session
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} LoginResponse "登录成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 401 {object} ErrorResponse "用户名或密码错误"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数格式错误", err)
		return
	}

	// 登录
	user, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) || errors.Is(err, services.ErrInvalidPassword) {
			Unauthorized(c, "用户名或密码错误", err)
			return
		} else {
			InternalError(c, "登录失败", err)
			return
		}
	}

	// session处理
	session, err := h.sessionStore.Get(c.Request, "library-session")
	if err != nil {
		InternalError(c, "Session初始化失败", err)
		return
	}

	session.Values["authenticated"] = true
	session.Values["userID"] = user.ID
	session.Values["username"] = user.Name
	session.Values["role"] = user.Role

	// 保存Session
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		InternalError(c, "无法保存Session", err)
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Message: "登录成功",
		User:    user,
	})
}

// Register godoc
// @Summary 用户注册
// @Description 新用户注册账号
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册信息"
// @Success 201 {object} RegisterResponse "注册成功"
// @Failure 400 {object} ErrorResponse "请求参数错误或用户名已存在"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数格式错误", err)
		return
	}

	// 注册用户
	err := h.authService.Register(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrUserExists) {
			BadRequest(c, "用户名已存在", err)
			return
		} else {
			InternalError(c, "用户注册失败", err)
			return
		}
	}

	c.JSON(http.StatusCreated, RegisterResponse{Message: "用户注册成功"})
}

// Logout godoc
// @Summary 用户注销
// @Description 用户注销登录，清除session。需要用户已登录。
// @Tags auth
// @Accept json
// @Produce json
// @Success 204 "注销成功"
// @Failure 401 {object} ErrorResponse "未登录用户无法注销"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 获取Session
	session, err := h.sessionStore.Get(c.Request, "library-session")
	if err != nil {
		InternalError(c, "Session错误", err)
		return
	}

	// 检查是否已登录
	if session.IsNew {
		Unauthorized(c, "未登录用户无法注销", nil)
		return
	}

	// 注销用户
	session.Options.MaxAge = -1
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		InternalError(c, "注销失败", err)
		return
	}

	c.Status(http.StatusNoContent)
}

// 请求和响应结构体定义
type RegisterRequest struct {
	Username string `json:"username" binding:"required" example:"user123"`
	Password string `json:"password" binding:"required" example:"password123"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"user123"`
	Password string `json:"password" binding:"required" example:"password123"`
}

type LoginResponse struct {
	Message string       `json:"message" example:"登录成功"`
	User    *models.User `json:"user"`
}

type RegisterResponse struct {
	Message string `json:"message" example:"用户注册成功"`
}
