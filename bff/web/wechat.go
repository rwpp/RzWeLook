package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	oauth2v1 "github.com/rwpp/RzWeLook/api/proto/gen/oauth2/v1"
	userv1 "github.com/rwpp/RzWeLook/api/proto/gen/user/v1"
	ijwt "github.com/rwpp/RzWeLook/bff/web/jwt"
	"net/http"
	"time"
)

type OAuthWechatHandler struct {
	svc     oauth2v1.Oauth2ServiceClient
	userSvc userv1.UserServiceClient
	ijwt.Handler
	stateKey []byte
	cfg      WechatHandlerConfig
}

type WechatHandlerConfig struct {
	Secure bool
}

func NewOAuthWechatHandler(svc oauth2v1.Oauth2ServiceClient,
	userSvc userv1.UserServiceClient,
	jwtHdl ijwt.Handler,
	cfg WechatHandlerConfig) *OAuthWechatHandler {
	return &OAuthWechatHandler{
		svc:      svc,
		stateKey: []byte("secret"),
		userSvc:  userSvc,
		cfg:      cfg,
		Handler:  jwtHdl,
	}
}

func (h *OAuthWechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth/wechat")
	g.GET("/authurl", h.AuthURL)
	g.Any("/callback", h.Callback)

}

func (h *OAuthWechatHandler) AuthURL(ctx *gin.Context) {
	state := uuid.New()
	url, err := h.svc.AuthURL(ctx, &oauth2v1.AuthURLRequest{
		State: state,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "生成URL失败",
		})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{
			//过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 3)),
		},
	})
	tokenStr, err := token.SignedString(h.stateKey)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "错误",
		})
		return
	}
	ctx.SetCookie("jwt-state", tokenStr, 600,
		"/oauth/wechat/callback", "", h.cfg.Secure, true)
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})
}

func (h *OAuthWechatHandler) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	err := h.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "登陆失败",
		})
		return
	}
	info, err := h.svc.VerifyCode(ctx, &oauth2v1.VerifyCodeRequest{
		Code: code,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	u, err := h.userSvc.FindOrCreateByWechat(ctx, &userv1.FindOrCreateByWechatRequest{
		Info: &userv1.WechatInfo{
			OpenId:  info.OpenId,
			UnionId: info.UnionId,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}

	err = h.SetLoginToken(ctx, u.User.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	ctx.JSON(http.StatusOK, Result{
		Data: u,
		Msg:  "成功",
	})
}

func (h *OAuthWechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	ck, err := ctx.Cookie("jwt-state")
	if err != nil {
		return fmt.Errorf("拿不到cookie: %w", err)
	}
	var sc StateClaims
	token, err := jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return h.stateKey, nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("token过期%w", err)
	}
	if sc.State != state {
		return errors.New("state不相等")
	}
	return nil
}

type StateClaims struct {
	State string
	jwt.RegisteredClaims
}
