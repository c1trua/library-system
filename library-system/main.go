package main

import (
	"fmt"
	"library-system/handlers"
	"library-system/middleware"
	"library-system/models"
	"library-system/repositories"
	"library-system/services"
	"log"
	"os"

	_ "library-system/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// @title 图书管理系统 API
// @version 1.0
// @description hhahaaaahahahahaha

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in cookie
// @name library-session
// @description 用户登录后，Session Cookie会自动携带在请求中
func main() {
	dbHost := getEnv("MYSQL_HOST", "localhost")
	dbUser := getEnv("MYSQL_USER", "root")
	dbPassword := getEnv("MYSQL_PASSWORD", "")
	dbName := getEnv("MYSQL_DBNAME", "library-system")
	sessionSecret := getEnv("SESSION_SECRET", "SBSBSBSBSBSSBSBS")
	serverPort := getEnv("SERVER_PORT", ":8080")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}
	err = db.AutoMigrate(&models.Book{})
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}
	err = db.AutoMigrate(&models.BorrowRecord{})
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	// 初始化Session
	sessionStore := sessions.NewCookieStore([]byte(sessionSecret))

	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   24 * 60 * 60,
		HttpOnly: true,
		Secure:   false,
	}

	// 初始化各层组件
	userRepo := repositories.NewUserRepository(db)
	bookRepo := repositories.NewBookRepository(db)
	authService := services.NewAuthService(userRepo)
	bookService := services.NewBookService(bookRepo)
	borrowService := services.NewBorrowService(db)
	adminService := services.NewAdminService(db)
	authHandler := handlers.NewAuthHandler(authService, sessionStore)
	bookHandler := handlers.NewBookHandler(bookService)
	borrowHandler := handlers.NewBorrowHandler(borrowService)
	adminHandler := handlers.NewAdminHandler(adminService)

	// 创建路由
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		// 认证路由
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
		}

		// 需要认证的路由
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(sessionStore))
		{
			// 图书路由
			books := protected.Group("/books")
			{
				books.GET("", bookHandler.GetAllBooks)                            // GET /api/v1/books
				books.GET("/search", bookHandler.SearchBooksByKeyword)            // GET /api/v1/books/search?keyword=xxx
				books.GET("/search/title", bookHandler.SearchBooksByTitleKeyword) // GET /api/v1/books/search/title?titlekeyword=xxx
				books.GET("/search/author", bookHandler.SearchBooksByAuthor)      // GET /api/v1/books/search/author?author=xxx
				books.GET("/:id", bookHandler.GetBookInfoByID)                    // GET /api/v1/books/1
				books.GET("/title", bookHandler.GetBookInfoByTitle)               // GET /api/v1/books/title?title=具体的标题
			}

			// 借阅路由
			borrow := protected.Group("/borrow")
			{
				borrow.POST("", borrowHandler.BorrowBook)                  // POST /api/v1/borrow
				borrow.POST("/return", borrowHandler.ReturnBook)           // POST /api/v1/borrow/return
				borrow.GET("/records", borrowHandler.GetUserBorrowRecords) // GET /api/v1/borrow/records
			}

			// 管理员路由
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.POST("/books", adminHandler.AddBook)                     // POST /api/v1/admin/books
				admin.PUT("/books", adminHandler.UpdateBook)                   // PUT /api/v1/admin/books
				admin.DELETE("/books", adminHandler.DeleteBook)                // DELETE /api/v1/admin/books
				admin.GET("/borrow-records", adminHandler.GetAllBorrowRecords) // GET /api/v1/admin/borrow-records
			}
		}
	}

	// 启动服务器
	log.Println("服务器启动")
	if err := router.Run(serverPort); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
