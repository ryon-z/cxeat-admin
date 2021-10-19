package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"yelloment-admin/common"
	"yelloment-admin/models"

	"github.com/gin-gonic/gin"
)

// ShowCodeListPage 코드 리스트 페이지
func ShowCodeListPage(c *gin.Context) {
	list := models.GetCodeList("", false)

	render(c, gin.H{
		"title": "코드관리",
		"list":  list,
	}, "code/list.html")
}

// ShowCodeViewPage 코드 상세 페이지
func ShowCodeViewPage(c *gin.Context) {
	codeNo, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Error(err)

		render(c, gin.H{
			"title": "Not Found",
		}, "error/notfound.html")

		return
	}

	typelist := models.GetCodeTypeList()
	info := models.GetCodeInfo(codeNo)

	render(c, gin.H{
		"title":    "코드상세",
		"typelist": typelist,
		"info":     info,
	}, "code/view.html")
}

// ShowCodeViewModalPage 코드 상세 페이지(모달)
func ShowCodeViewModalPage(c *gin.Context) {
	codeNo, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Error(err)

		render(c, gin.H{
			"title": "Not Found",
		}, "error/notfound.html")

		return
	}

	typelist := models.GetCodeTypeList()
	info := models.GetCodeInfo(codeNo)

	render(c, gin.H{
		"title":    "코드상세",
		"typelist": typelist,
		"info":     info,
	}, "code/view_modal.html")
}

// SetCode 코드 등록/수정
func SetCode(c *gin.Context) {
	var info models.CodeMst

	err := c.ShouldBind(&info)

	if err != nil {
		log.Println(err)
	}

	models.SetCode(info)

	go common.SendSlackMessage("system", fmt.Sprintf("코드변경 - %s", GetAdminID(c)))

	c.JSON(http.StatusOK, info)
}
