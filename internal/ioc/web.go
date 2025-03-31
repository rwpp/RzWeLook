package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rwpp/RzWeLook/internal/web"
	"github.com/rwpp/RzWeLook/internal/web/middleware"
	"github.com/rwpp/RzWeLook/pkg/ginx/middlewares/ratelimit"
	"strings"
	"time"
)

func InitWeb(mdls []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	// 这里可以初始化web服务
	// 比如使用gin框架
	// 实际上什么都不做
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	return server
}
func InitMiddleware(redisClient redis.Cmdable) []gin.HandlerFunc {
	// 这里可以初始化中间件
	// 比如使用gin框架
	// 实际上什么都不做
	return []gin.HandlerFunc{
		corsHdl(),
		middleware.NewLoginJWTMiddlewareBuilder().Build(),
		ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
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
	})
}
