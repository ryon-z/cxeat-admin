package controllers

import (
	"github.com/gin-gonic/gin"
)

// ShowDashboardPage 대시보드 페이지
func ShowDashboardPage(c *gin.Context) {
	render(c, gin.H{
		"title": "대시보드",
	}, "dashboard/main.html")
}
