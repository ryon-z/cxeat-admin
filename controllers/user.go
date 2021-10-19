package controllers

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"yelloment-admin/models"

	"github.com/gin-gonic/gin"
)

// ShowUserListPage 회원 리스트 페이지
func ShowUserListPage(c *gin.Context) {
	render(c, gin.H{
		"title": "회원관리",
	}, "user/list.html")
}

// ShowUserListPage 회원 리스트 페이지
func ShowLeaveUserListPage(c *gin.Context) {
	render(c, gin.H{
		"title": "탈퇴회원관리",
	}, "user/leavelist.html")
}

// ShowUserViewPage 회원상세 페이지
func ShowUserViewPage(c *gin.Context) {
	userno, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.Error(err)

		render(c, gin.H{
			"title": "Not Found",
		}, "error/notfound.html")

		return
	}

	info := models.GetUserInfo(userno)

	if info.UserNo < 1 {
		render(c, gin.H{
			"title": "Not Found",
		}, "error/notfound.html")

		return
	}

	addresslist := models.GetUserAddressList(info.UserNo)
	cardlist := models.GetUserCardList(info.UserNo)

	render(c, gin.H{
		"title":       "회원상세",
		"info":        info,
		"addresslist": addresslist,
		"cardlist":    cardlist,
	}, "user/view.html")
}

// ShowUserViewPopPage 회원상세 페이지(팝업)
func ShowUserViewPopPage(c *gin.Context) {
	render(c, gin.H{
		"title": "회원상세",
	}, "user/view_popup.html")
}

// GetUserList 회원 조회
func GetUserList(c *gin.Context) {
	status := strings.Split(c.DefaultQuery("status", ""), "|")
	list := models.GetUserList(status)

	c.JSON(http.StatusOK, gin.H{"data": list})
}

// GetUserAddressList 회원 주소록 조회
func GetUserAddressList(c *gin.Context) {
	uNo, err := strconv.ParseInt(c.DefaultQuery("userno", "0"), 10, 64)
	if err != nil {
		log.Println(err)

		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	list := models.GetUserAddressList(uNo)

	c.JSON(http.StatusOK, gin.H{"data": list})
}

func GetUserAddress(c *gin.Context) {
	aNo, err := strconv.ParseInt(c.DefaultQuery("addressno", "0"), 10, 64)
	if err != nil {
		log.Println(err)

		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	info := models.GetUserAddress(aNo)
	if info.AddressNo < 1 {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": info})
}
