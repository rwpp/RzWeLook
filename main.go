package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/rwpp/RzWeLook/config"
	"github.com/rwpp/RzWeLook/internal/repository/cache"
	"github.com/rwpp/RzWeLook/internal/service/sms/memory"

	//"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rwpp/RzWeLook/internal/repository"
	"github.com/rwpp/RzWeLook/internal/repository/dao"
	"github.com/rwpp/RzWeLook/internal/service"
	"github.com/rwpp/RzWeLook/internal/web"
	"github.com/rwpp/RzWeLook/internal/web/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

func main() {
	connStr := "root:root@tcp(localhost:13316)/welook"
	db, err := gorm.Open(mysql.Open(connStr))
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	// 创建新表
	err = dao.InitTable(db)
	ud := dao.NewUserDAO(db)
	uc := cache.NewUserCache(redisClient)
	repo := repository.NewUserRepository(ud, uc)
	svc := service.NewUserService(repo)
	codeCache := cache.NewCodeCache(redisClient)
	codeRepo := repository.NewCodeRepository(codeCache)
	smsSvc := memory.NewService()
	codeSvc := service.NewCodeService(codeRepo, smsSvc)
	u := web.NewUserHandler(svc, codeSvc)
	server := gin.Default()
	//server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
	server.Use(cors.New(cors.Config{
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"X-Jwt-Token", "x-refresh-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	}))
	store := memstore.NewStore([]byte("NtgEPQxuMoH3aCLuQW2NaAy3FoL3tveW"),
		[]byte("soiccWYEPOhb5sR4LLs6OnUT3Gf45JL1"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("mysession", store))
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().Build())
	u.RegisterRoutes(server)
	log.Println("Server starting on :8080")
	server.Run(":8080")
}
func initDB() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		return nil, err
	}
	return db, nil
}
