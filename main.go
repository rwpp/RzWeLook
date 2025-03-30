package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/rwpp/RzWeLook/internal/repository"
	"github.com/rwpp/RzWeLook/internal/repository/dao"
	"github.com/rwpp/RzWeLook/internal/service"
	"github.com/rwpp/RzWeLook/internal/web"
	"github.com/rwpp/RzWeLook/internal/web/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {
	log.Println("Starting application...")
	connStr := "root:root@tcp(localhost:13316)/welook"
	db, err := gorm.Open(mysql.Open(connStr))
	if err != nil {

		panic(err)
	}
	// 创建新表
	err = dao.InitTable(db)
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	server := gin.Default()
	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8080"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.ExposeHeaders = []string{"x-jwt-token"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	config.AllowCredentials = true
	server.Use(cors.New(config))
	// 注册路由
	//store := cookie.NewStore([]byte("secret"))
	//store := memstore.NewStore([]byte("soiccWYEPOhb5sR4LLs6OnUT3Gf45JL1"),
	//[]byte("NtgEPQxuMoH3aCLuQW2NaAy3FoL3tveW"))

	store, err := redis.NewStore(16,
		"tcp", "localhost:6379", "",
		[]byte("NtgEPQxuMoH3aCLuQW2NaAy3FoL3tveW"),
		[]byte("soiccWYEPOhb5sR4LLs6OnUT3Gf45JL1"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("mysession", store))

	//server.Use(middleware.NewLoginMiddlewareBuilder().
	//	//IgnorePaths("/users/login").
	//	//IgnorePaths("users/signup").
	//	Build())
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/login").
		IgnorePaths("/users/signup").
		Build())
	u.RegisterRoutes(server)
	log.Println("Server starting on :8080")
	server.Run(":8080")
}
