package ioc

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rwpp/RzWeLook/bff/web"
	ijwt "github.com/rwpp/RzWeLook/bff/web/jwt"
	"github.com/rwpp/RzWeLook/bff/web/middleware"
	"github.com/rwpp/RzWeLook/pkg/ginx"
	"github.com/rwpp/RzWeLook/pkg/logger"
	"github.com/spf13/viper"
	"strings"
	"time"
)

func InitGinServer(
	l logger.LoggerV1,
	jwtHdl ijwt.Handler,
	user *web.UserHandler,
	article *web.ArticleHandler,
	reward *web.RewardHandler,
) *ginx.Server {
	engine := gin.Default()
	engine.Use(
		corsHdl(),
		timeout(),
		middleware.NewLoginJWTMiddlewareBuilder(jwtHdl).Build())
	user.RegisterRoutes(engine)
	article.RegisterRoutes(engine)
	reward.RegisterRoutes(engine)
	addr := viper.GetString("http.addr")
	ginx.InitCounter(prometheus.CounterOpts{
		Namespace: "daming_geektime",
		Subsystem: "webook_bff",
		Name:      "http",
	})
	ginx.SetLogger(l)
	return &ginx.Server{
		Engine: engine,
		Addr:   addr,
	}
}

func timeout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, ok := ctx.Request.Context().Deadline()
		if !ok {
			// 强制给一个超时，省得我前端调试等得不耐烦
			newCtx, cancel := context.WithTimeout(ctx.Request.Context(), time.Second*10)
			defer cancel()
			ctx.Request = ctx.Request.Clone(newCtx)
		}
		ctx.Next()
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins: []string{"*"},
		//AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 你不加这个，前端是拿不到的
		ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
		// 是否允许你带 cookie 之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 你的开发环境
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
