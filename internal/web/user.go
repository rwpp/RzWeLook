package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/rwpp/RzWeLook/internal/domain"
	"github.com/rwpp/RzWeLook/internal/service"
	"net/http"
)

type UserHandler struct {
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	svc         service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	const (
		emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		// 和上面比起来，用 ` 看起来就比较清爽
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`

		userIdKey = "userId"
		bizLogin  = "login"
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		emailExp:    emailExp,
		passwordExp: passwordExp,
		svc:         svc,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.GET("/edit", u.Edit)
	ug.GET("/logout", u.Logout)
	ug.GET("/profile", u.Profile)
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
	if err == service.ErrUserDuplicateEmail {

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
	err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "账号或密码错误")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) Profile(ctx *gin.Context) {

}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Logout(ctx *gin.Context) {

}
