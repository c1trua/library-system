package middleware

import (
	"library-system/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// AuthMiddleware
func AuthMiddleware(sessionStore sessions.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Session
		session, err := sessionStore.Get(c.Request, "library-session")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "无法获取Session"})
			return
		}
		// 检查是否已认证
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未授权，请先登录"})
			return
		}
		// 从Session中获取用户信息，并存入Gincontext
		userID := session.Values["userID"].(int)
		username := session.Values["username"].(string)
		role := session.Values["role"].(string)

		user := &models.User{
			ID:   userID,
			Name: username,
			Role: role,
		}

		c.Set("user", user)
		c.Next()
	}
}

// AdminMiddleware
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Gincontext获取用户信息
		userObj, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "用户信息不存在"})
			return
		}

		user := userObj.(*models.User)
		// 检查用户角色是否为管理员
		if user.Role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "权限不足，需要管理员权限"})
			return
		}
		c.Next()
	}
}
