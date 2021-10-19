package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"yelloment-admin/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var RandNo int

func init() {
	RandNo = rand.Int()
}

func render(c *gin.Context, data gin.H, templateName string) {
	session := sessions.Default(c)

	if aid := session.Get("adminid"); aid != nil {
		data["admininfo"] = models.GetAdminAccount(aid.(string))
	}

	data["reqURL"] = c.Request.RequestURI
	data["randNo"] = RandNo

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		// Respond with XML
		c.XML(http.StatusOK, data["payload"])
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, templateName, data)
	}
}

func GetAdminID(c *gin.Context) string {
	return sessions.Default(c).Get("adminid").(string)
}

// FileUpload 파일 업로드
func FileUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	//임시로 웹서버에 업로드 하는 코드로 작성
	//추후 aws s3로 업로드 하는 로직으로 변경 예정
	filename := "contents/upload/" + file.Filename

	fmt.Println(filename)

	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully.", file.Filename))
}
