package web

import "github.com/gin-gonic/gin"

type UserHandler struct {
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/user")
	ug.GET("/profile", u.Profile)
	ug.GET("/signup", u.SignUp)
	ug.GET("/login", u.Login)
	ug.GET("/edit", u.Edit)
	ug.GET("/logout", u.Logout)

}

func (u *UserHandler) Profile(context *gin.Context) {

}

func (u *UserHandler) SignUp(context *gin.Context) {

}

func (u *UserHandler) Login(context *gin.Context) {

}

func (u *UserHandler) Edit(context *gin.Context) {

}

func (u *UserHandler) Logout(context *gin.Context) {

}
