package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rwpp/RzWeLook/pkg/logger"
	"net/http"
)

var L logger.LoggerV1

func WrapToken[C jwt.Claims](fn func(ctx *gin.Context, uc C) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		val, ok := ctx.Get("users")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c, ok := val.(C)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		res, err := fn(ctx, c)
		if err != nil {

			L.Error("业务处理逻辑错误",
				logger.String("path", ctx.Request.URL.Path),
				logger.String("route", ctx.FullPath()),

				logger.Error(err))

		}
		ctx.JSON(http.StatusOK, res)
	}
}
func WrapBodyAndToken[Req any, C jwt.Claims](fn func(ctx *gin.Context, req Req, uc C) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			return
		}

		var c C
		val, ok := ctx.Get("users")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c, ok = val.(C)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		res, err := fn(ctx, req, c)
		if err != nil {

			L.Error("业务处理逻辑错误",
				logger.String("path", ctx.Request.URL.Path),
				logger.String("route", ctx.FullPath()),

				logger.Error(err))

		}
		ctx.JSON(http.StatusOK, res)
	}
}

func WrapBodyV1[T any](fn func(ctx *gin.Context, req T) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			return
		}
		res, err := fn(ctx, req)
		if err != nil {

			L.Error("业务处理逻辑错误",
				logger.String("path", ctx.Request.URL.Path),
				logger.String("route", ctx.FullPath()),

				logger.Error(err))

		}
		ctx.JSON(http.StatusOK, res)
	}
}

func WrapBody[T any](l logger.LoggerV1, fn func(ctx *gin.Context, req T) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			return
		}
		res, err := fn(ctx, req)
		if err != nil {

			l.Error("业务处理逻辑错误",
				logger.String("path", ctx.Request.URL.Path),
				logger.String("route", ctx.FullPath()),

				logger.Error(err))

		}
		ctx.JSON(http.StatusOK, res)
	}
}
