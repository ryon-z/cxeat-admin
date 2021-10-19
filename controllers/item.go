package controllers

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"yelloment-admin/models"

	"github.com/gin-gonic/gin"
)

// ShowItemListPage 상품 리스트 페이지
func ShowItemListPage(c *gin.Context) {
	render(c, gin.H{
		"title": "상품관리",
	}, "item/list.html")
}

// ShowItemViewPage 상품 상세 페이지
func ShowItemViewPage(c *gin.Context) {
	itemNo, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Error(err)

		render(c, gin.H{
			"title": "Not Found",
		}, "error/notfound.html")

		return
	}

	catelist := models.GetCodeList("ITEM_CATEGORY", false)
	unitlist := models.GetCodeList("ITEM_UNIT_TYPE", false)
	info := models.GetItemInfo(itemNo)

	render(c, gin.H{
		"title":    "상품상세",
		"catelist": catelist,
		"unitlist": unitlist,
		"info":     info,
	}, "item/view.html")
}

// SetItem 상품 등록/수정
func SetItem(c *gin.Context) {
	var info models.ItemMst

	err := c.ShouldBind(&info)

	if err != nil {
		log.Println(err)
	}

	models.SetItem(info)

	c.JSON(http.StatusOK, info)
}

func GetItemList(c *gin.Context) {
	status := strings.Split(c.DefaultQuery("status", ""), "|")

	list := models.GetItemList(status)

	c.JSON(http.StatusOK, gin.H{"data": list})
}

// ShowBundleListPage 번들 리스트 페이지
func ShowBundleListPage(c *gin.Context) {
	render(c, gin.H{
		"title": "상품 번들 관리",
	}, "item/bundlelist.html")
}

// ShowBundleViewPage 번들 상세 페이지
func ShowBundleViewPage(c *gin.Context) {
	bundleNo, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Error(err)

		render(c, gin.H{
			"title": "Not Found",
		}, "error/notfound.html")

		return
	}

	info := models.GetBundleInfo(bundleNo)

	render(c, gin.H{
		"title": "상품 번들 상세",
		"info":  info,
	}, "item/bundleview.html")
}

// GetBundleInfo 번들 상세 페이지
func GetBundleInfo(c *gin.Context) {
	bundleNo, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusNotFound, gin.H{})

		return
	}

	info := models.GetBundleInfo(bundleNo)

	c.JSON(http.StatusOK, gin.H{"data": info})
}

// BundleManage 구성 상품 벌크
type BundleManage struct {
	BundleNo   string
	BundleName string
	StatusCode string

	Items []ItemManage `json:"items"`
}

type ItemManage struct {
	ItemNo  string
	ItemCnt string
}

// SetBundle 번들 등록/수정
func SetBundle(c *gin.Context) {
	var info BundleManage

	err := c.ShouldBind(&info)

	if err != nil {
		log.Println(err)
	}

	bNo, _ := strconv.Atoi(info.BundleNo)

	result := models.SetBundle(models.BundleMst{BundleNo: bNo, BundleName: info.BundleName, StatusCode: info.StatusCode})

	bNo = int(result)

	models.InitBundleItem(bNo)

	for _, i := range info.Items {
		iNo, _ := strconv.Atoi(i.ItemNo)
		iCt, _ := strconv.Atoi(i.ItemCnt)

		models.AddBundleItem(models.BundleItem{BundleNo: bNo, ItemNo: iNo, ItemCnt: iCt})
	}

	c.JSON(http.StatusOK, gin.H{"id": bNo})
}

// RemoveBundle 번들 삭제
func RemoveBundle(c *gin.Context) {
	bundleNo, err := strconv.Atoi(c.DefaultPostForm("bundleno", "0"))
	if err != nil {
		log.Println(err)
	}

	result := models.RemoveBundle(bundleNo)
	if result < 1 {
		c.JSON(http.StatusNotFound, gin.H{})
	}

	c.JSON(http.StatusOK, gin.H{"id": bundleNo})
}

func GetBundleList(c *gin.Context) {
	status := strings.Split(c.DefaultQuery("status", ""), "|")

	list := models.GetBundleList(status)

	c.JSON(http.StatusOK, gin.H{"data": list})
}
