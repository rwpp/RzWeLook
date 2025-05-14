package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	codev1 "github.com/rwpp/RzWeLook/api/proto/gen/code/v1"
	userv1 "github.com/rwpp/RzWeLook/api/proto/gen/user/v1"
	ijwt "github.com/rwpp/RzWeLook/internal/web/jwt"
	"github.com/rwpp/RzWeLook/pkg/ginx"
	"github.com/rwpp/RzWeLook/pkg/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	svc         userv1.UserServiceClient
	codeSvc     codev1.CodeServiceClient
	ijwt.Handler
	cmd redis.Cmdable
	l   logger.LoggerV1
}

func NewUserHandler(svc userv1.UserServiceClient, codeSvc codev1.CodeServiceClient,
	jwtHdl ijwt.Handler,
	l logger.LoggerV1) *UserHandler {
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
		l:           l,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	//ug.POST("/login", u.Login)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/login_sms",
		ginx.WrapBody[LoginSMSReq](u.l.With(logger.String("method", "login_sms")), u.LoginSMS))
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

	_, err := u.codeSvc.Send(ctx, &codev1.CodeSendRequest{
		Biz: biz, Phone: req.Phone,
	})
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
		//TODO:用grpc传递错误码
	//case service.ErrCodeSendTooMany:
	//	ctx.JSON(http.StatusOK, Result{
	//		Code: 4,
	//		Msg:  "发送验证码过于频繁",
	//	})
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
	_, err = u.svc.Signup(ctx.Request.Context(),
		&userv1.SignupRequest{User: &userv1.User{Email: req.Email,
			Password: req.ConfirmPassword}})
	//if errors.Is(err, service.ErrUserDuplicate) {
	//	ctx.JSON(http.StatusOK, gin.H{"error": "邮箱已被注册"})
	//	return
	//}
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": "系统异常"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

//func (u *UserHandler) Login(ctx *gin.Context) {
//	type LoginReq struct {
//		Email    string `json:"email"`
//		Password string `json:"password"`
//	}
//	var req LoginReq
//	if err := ctx.Bind(&req); err != nil {
//		return
//	}
//	user, err := u.svc.Login(ctx, req.Email, req.Password)
//	if err == service.ErrInvalidUserOrPassword {
//		ctx.String(http.StatusOK, "账号或密码错误")
//		return
//	}
//	if err != nil {
//		ctx.String(http.StatusOK, "系统异常")
//		return
//	}
//	sess := sessions.Default(ctx)
//	sess.Set("userId", user.Id)
//	sess.Options(sessions.Options{
//		//Secure:   true,
//		//HttpOnly: true,
//		MaxAge: 60 * 30,
//	})
//	sess.Save()
//	ctx.String(http.StatusOK, "登录成功")
//	return
//}

//	func (u *UserHandler) Profile(ctx *gin.Context) {
//		type Profile struct {
//			Email string `json:"email"`
//		}
//		sess := sessions.Default(ctx)
//		id := sess.Get(userIdKey).(int64)
//		ux, err := u.svc.Profile(ctx, id)
//		if err != nil {
//			ctx.JSON(http.StatusOK, Result{
//				Code: 1,
//				Msg:  "系统错误",
//			})
//			return
//		}
//		ctx.JSON(http.StatusOK, Profile{
//			Email: ux.Email,
//		})
//
// }
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
	uc := ctx.MustGet("users").(ijwt.UserClaims)
	_, err = u.svc.UpdateNonSensitiveInfo(ctx, &userv1.UpdateNonSensitiveInfoRequest{
		User: &userv1.User{
			Id:       uc.Uid,
			Nickname: req.Nickname,
			AboutMe:  req.AboutMe,
			Birthday: timestamppb.New(birthday),
		},
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
	user, err := u.svc.Login(ctx, &userv1.LoginRequest{
		Email: req.Email, Password: req.Password,
	})
	//if err == service.ErrInvalidUserOrPassword {
	//	ctx.String(http.StatusOK, "账号或密码错误")
	//	return
	//}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	if err = u.SetLoginToken(ctx, user.User.Id); err != nil {
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
	//ToDo 可能key不对
	uc := ctx.MustGet("users").(ijwt.UserClaims)
	ux, err := u.svc.Profile(ctx, &userv1.ProfileRequest{Id: uc.Uid})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 1,
			Msg:  "系统错误",
		})
		return
	}
	m := ux.User
	ctx.JSON(http.StatusOK, Profile{
		Email:    m.Email,
		Phone:    m.Phone,
		Nickname: m.Nickname,
		Birthday: m.Birthday.AsTime().Format(time.DateOnly),
		AboutMe:  m.AboutMe,
	})

}

type LoginSMSReq struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

func (u *UserHandler) LoginSMS(ctx *gin.Context, req LoginSMSReq) (ginx.Result, error) {
	resp, err := u.codeSvc.Verify(ctx, &codev1.VerifyRequest{
		Biz: biz, Phone: req.Phone, InputCode: req.Code,
	})
	if err != nil {
		return Result{Code: 5, Msg: "系统异常"}, err
	}
	if resp.Answer {
		return Result{Code: 5, Msg: "验证码错误"}, nil
	}
	m, err := u.svc.FindOrCreate(ctx, &userv1.FindOrCreateRequest{
		Phone: req.Phone,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	//Todo: 这里需要根据手机号查询用户信息
	ssid := uuid.New().String()
	if err = u.SetJWToken(ctx, m.User.Id, ssid); err != nil {
		return Result{Code: 5, Msg: "系统异常"}, err
	}
	return Result{Msg: "登录成功"}, nil
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
