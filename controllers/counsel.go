package controllers

import (
	"log"
	"net/http"
	"strconv"
	"yelloment-admin/models"

	"github.com/gin-gonic/gin"
)

// GetCounselList 상담내역 조회
func GetCounselList(c *gin.Context) {
	uNo, _ := strconv.ParseInt(c.DefaultQuery("userno", "0"), 10, 64)
	sNo, _ := strconv.ParseInt(c.DefaultQuery("subsno", "0"), 10, 64)
	oNo, _ := strconv.ParseInt(c.DefaultQuery("orderno", "0"), 10, 64)

	list := models.GetCounselList(uNo, sNo, oNo)

	c.JSON(http.StatusOK, gin.H{"data": list})
}

// AddCounsel 상담 입력
func AddCounsel(c *gin.Context) {
	var info models.CounselMst

	err := c.ShouldBind(&info)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	aid := GetAdminID(c)
	info.RegUser = &aid

	result := models.AddCounsel(info)

	c.JSON(http.StatusOK, gin.H{"id": result})
}

func RemoveCounsel(c *gin.Context) {
	cNo, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println(err)

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	aid := GetAdminID(c)

	result := models.RemoveCounsel(cNo, aid)

	c.JSON(http.StatusOK, gin.H{"id": result})
}
