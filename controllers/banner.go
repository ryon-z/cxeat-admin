package controllers

import (
	"log"
	"net/http"
	"strconv"
	"yelloment-admin/models"

	"github.com/gin-gonic/gin"
)

// ShowBannerListPage 배너 리스트 페이지
func ShowBannerListPage(c *gin.Context) {
	list := models.GetBannerList()

	render(c, gin.H{
		"title": "배너관리",
		"list":  list,
	}, "banner/list.html")
}

// ShowBannerViewPage 배너 상세 페이지
func ShowBannerViewPage(c *gin.Context) {
	bannerno, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Error(err)

		render(c, gin.H{
			"title": "Not Found",
		}, "error/notfound.html")

		return
	}

	info := models.GetBannerInfo(bannerno)

	if info.BannerNo < 1 {
		render(c, gin.H{
			"title": "배너상세",
			"info":  info,
		}, "banner/view.html")

		return
	}

	render(c, gin.H{
		"title": "배너상세",
		"info":  info,
	}, "banner/view.html")
}

// SetBanner 배너 등록/수정
func SetBanner(c *gin.Context) {
	var info models.BannerMst

	err := c.ShouldBind(&info)

	if err != nil {
		log.Println(err)
	}

	models.SetBanner(info)

	c.JSON(http.StatusOK, info)
}
