package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler interface {
	SetLoginToken(ctx *gin.Context, uid int64) error
	SetJWToken(ctx *gin.Context, uid int64, ssid string) error
	ClearToken(ctx *gin.Context) error
	CheckSession(ctx *gin.Context, ssid string) error
	ExtractUser(ctx *gin.Context) string
}
type RefreshClaims struct {
	Uid  int64
	Ssid string
	jwt.RegisteredClaims
}
type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	Ssid      string
	UserAgent string
}
