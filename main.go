package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rwpp/RzWeLook/internal/repository"
	"github.com/rwpp/RzWeLook/internal/repository/dao"
	"github.com/rwpp/RzWeLook/internal/service"
	"github.com/rwpp/RzWeLook/internal/web"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {
	log.Println("Starting application...")
	connStr := "root:root@tcp(localhost:13316)/welook"
	db, err := gorm.Open(mysql.Open(connStr))
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		panic(err)
	}
	log.Println("Database connected successfully")

	// 删除旧表
	err = db.Migrator().DropTable(&dao.User{})
	if err != nil {
		log.Printf("Failed to drop table: %v", err)
	}
	log.Println("Old tables dropped successfully")

	// 创建新表
	err = dao.InitTable(db)
	if err != nil {
		log.Printf("Failed to initialize tables: %v", err)
		panic(err)
	}
	log.Println("Tables initialized successfully")

	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	log.Println("User repository created successfully")
	svc := service.NewUserService(repo)
	log.Println("User service created successfully")
	u := web.NewUserHandler(svc)
	log.Println("User handler created successfully")

	server := gin.Default()

	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8080"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	config.AllowCredentials = true
	server.Use(cors.New(config))

	// 注册路由
	u.RegisterRoutes(server)

	log.Println("Server starting on :8080")
	server.Run(":8080")
}
