package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	ijwt "github.com/rwpp/RzWeLook/internal/web/jwt"
	"net/http"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
	ijwt.Handler
}

var IgnorePaths []string

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}
func NewLoginJWTMiddlewareBuilder(jwtHdl ijwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: jwtHdl,
	}
}
func (j *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不需要校验
		if ctx.Request.URL.Path == "/users/signup" ||
			ctx.Request.URL.Path == "/users/login" ||
			ctx.Request.URL.Path == "/users/refresh_token" ||
			ctx.Request.URL.Path == "/users/login_sms" ||
			ctx.Request.URL.Path == "/users/login_sms/code/send" ||
			ctx.Request.URL.Path == "/oauth/wechat/authurl" ||
			ctx.Request.URL.Path == "/oauth/wechat/callback" {
			return
		}
		// Authorization 头部
		// 得到的格式 Bearer token
		tokenStr := j.ExtractUser(ctx)
		claims := &ijwt.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("NtgEPQxuMoH3aCLuQW2NaAy3FoL3tveW"), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid || claims.Uid == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		err = j.CheckSession(ctx, claims.Ssid)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("claims", claims)
	}
}
