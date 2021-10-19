package controllers

import (
	"net/http"
	"yelloment-admin/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// ShowLoginPage 로그인 화면
func ShowLoginPage(c *gin.Context) {
	render(c, gin.H{
		"title": "로그인",
	}, "account/login.html")
}

// ConfirmLogin 로그인 확인
func ConfirmLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	returnurl, _ := c.GetQuery("returnurl")

	if models.IsAdminValid(username, password) {
		session := sessions.Default(c)
		session.Set("isauth", true)
		session.Set("adminid", username)
		session.Options(sessions.Options{
			Path:   "/",
			MaxAge: 3600 * 6, //세션 만료 시간 설정
		})
		session.Save()

		c.Redirect(http.StatusFound, returnurl)

	} else {
		render(c, gin.H{
			"title":        "로그인",
			"ErrorTitle":   "로그인 실패",
			"ErrorMessage": "ID 와 Password를 다시 확인해주세요."}, "account/login.html")
	}
}

// Logout 로그아웃
func Logout(c *gin.Context) {
	sessions.Default(c).Clear()
	sessions.Default(c).Save()

	c.Redirect(http.StatusFound, "/")
}
