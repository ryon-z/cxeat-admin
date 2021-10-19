package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// ConfirmAuth 인증 확인
func ConfirmAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		isauth := session.Get("isauth")
		if isauth == nil || isauth == "" || isauth == "false" {
			c.Redirect(http.StatusFound, "/auth/login?returnurl="+c.Request.RequestURI)
			c.Abort()
			return
		}

		adminid := session.Get("adminid")
		if adminid == nil || adminid == "" {
			c.Redirect(http.StatusFound, "/auth/login?returnurl="+c.Request.RequestURI)
			c.Abort()
			return
		}

		c.Next()
	}
}

// ConfirmAuthApi 인증 확인 API용
func ConfirmAuthApi() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		isauth := session.Get("isauth")
		if isauth == nil || isauth == "" || isauth == "false" {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		adminid := session.Get("adminid")
		if adminid == nil || adminid == "" {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Next()
	}
}
