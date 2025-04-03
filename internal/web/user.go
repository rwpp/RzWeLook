package web

import (
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/rwpp/RzWeLook/internal/domain"
	"github.com/rwpp/RzWeLook/internal/service"
	ijwt "github.com/rwpp/RzWeLook/internal/web/jwt"
	"log"
	"net/http"
	"time"
)

const biz = "login"
const userIdKey = "userId"

var _ handler = &UserHandler{}

type UserHandler struct {
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	svc         service.UserService
	codeSvc     service.CodeServiceInterface
	ijwt.Handler
	cmd redis.Cmdable
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeServiceInterface,
	jwtHdl ijwt.Handler) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	emailExp, err := regexp.Compile(emailRegexPattern, regexp.None)
	if err != nil {
		log.Printf("Failed to compile email regex: %v", err)
		panic("failed to compile email regex")
	}

	passwordExp, err := regexp.Compile(passwordRegexPattern, regexp.None)
	if err != nil {
		log.Printf("Failed to compile password regex: %v", err)
		panic("failed to compile password regex")
	}

	return &UserHandler{
		emailExp:    emailExp,
		passwordExp: passwordExp,
		svc:         svc,
		codeSvc:     codeSvc,
		Handler:     jwtHdl,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	//ug.POST("/login", u.Login)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/login_sms", u.LoginSMS)
	ug.POST("/login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("/edit", u.Edit)
	ug.POST("/logout", u.LogoutJWT)
	ug.GET("/profile", u.ProfileJWT)
	ug.GET("/refresh_token", u.RefreshToken)
}

func (u *UserHandler) RefreshToken(ctx *gin.Context) {

	refreshToken := u.ExtractUser(ctx)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(refreshToken, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.RtKey, nil
	})
	if err != nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = u.CheckSession(ctx, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
	err = u.SetJWToken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "刷新成功",
	})
}
func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type SendLoginSMSCodeReq struct {
		Phone string `json:"phone"`
	}
	var req SendLoginSMSCodeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	err := u.codeSvc.Send(ctx, biz, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "发送验证码过于频繁",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	return
}
func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}
	var req SignUpReq
	if err := ctx.ShouldBindJSON(&req); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": "系统错误"})
		return
	}
	if !ok {

		ctx.JSON(http.StatusOK, gin.H{"error": "邮箱格式错误"})
		return
	}
	if req.ConfirmPassword != req.Password {
		ctx.JSON(http.StatusOK, gin.H{"error": "两次密码不一致"})
		return
	}

	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": "系统错误"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"error": "密码格式错误"})
		return
	}
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserDuplicate) {
		ctx.JSON(http.StatusOK, gin.H{"error": "邮箱已被注册"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": "系统异常"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "账号或密码错误")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	sess := sessions.Default(ctx)
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		//Secure:   true,
		//HttpOnly: true,
		MaxAge: 60 * 30,
	})
	sess.Save()
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	type Profile struct {
		Email string `json:"email"`
	}
	sess := sessions.Default(ctx)
	id := sess.Get(userIdKey).(int64)
	ux, err := u.svc.Profile(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 1,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Profile{
		Email: ux.Email,
	})

}
func (u *UserHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"about_me"`
	}
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 1,
			Msg:  "请求参数错误",
		})
	}
	if req.Nickname == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "昵称不能为空",
		})
		return
	}
	if len(req.AboutMe) > 1024 {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "关于我的信息过长",
		})
		return
	}
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 6,
			Msg:  "生日格式错误",
		})
		return
	}
	uc := ctx.MustGet("user").(*ijwt.UserClaims)
	err = u.svc.UpdateNonSensitiveInfo(ctx, domain.User{
		Id:       uc.Uid,
		Nickname: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "修改成功",
	})
}
func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		MaxAge: -1})
	sess.Save()
	ctx.String(http.StatusOK, "注销成功")
}
func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "账号或密码错误")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	if err = u.SetLoginToken(ctx, user.Id); err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	fmt.Println(user)
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	type Profile struct {
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"about_me"`
	}

	uc := ctx.MustGet("user").(*ijwt.UserClaims)
	ux, err := u.svc.Profile(ctx, uc.Uid)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 1,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Profile{
		Email:    ux.Email,
		Phone:    ux.Phone,
		Nickname: ux.Nickname,
		Birthday: ux.Birthday.String(),
		AboutMe:  ux.AboutMe,
	})

}
func (u *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}
	ok, err := u.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码错误",
		})
		return
	}
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	//Todo: 这里需要根据手机号查询用户信息

	if err = u.SetLoginToken(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "验证码校验通过",
	})
}

func (u *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := u.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "退出登陆失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "退出登陆成功",
	})
}
