package controllers

import (
	"net/http"
	"strconv"
	"yelloment-admin/models"

	"github.com/gin-gonic/gin"
)

// GetOrderReview 주문 리뷰 조회
func GetReviewList(c *gin.Context) {
	userNo, _ := strconv.ParseInt(c.DefaultQuery("userno", "0"), 10, 64)
	subsNo, _ := strconv.ParseInt(c.DefaultQuery("subsno", "0"), 10, 64)
	orderNo, _ := strconv.ParseInt(c.DefaultQuery("orderno", "0"), 10, 64)

	list := models.GetReviewList(userNo, subsNo, orderNo)

	c.JSON(http.StatusOK, gin.H{"data": list})
}
