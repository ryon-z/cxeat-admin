package controllers

import (
	"net/http"
	"strconv"
	"yelloment-admin/models"

	"github.com/gin-gonic/gin"
)

// GetOrderList 주문 조회
func GetTagInfo(c *gin.Context) {
	tgNo, _ := strconv.ParseInt(c.DefaultQuery("taggroupno", "0"), 10, 64)

	list := models.GetTagInfo(tgNo)

	c.JSON(http.StatusOK, gin.H{"data": list})
}
