package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"yelloment-admin/common"
	"yelloment-admin/models"

	"github.com/gin-gonic/gin"
)

// ShowSubsListPage 상품 리스트 페이지
func ShowSubsListPage(c *gin.Context) {
	render(c, gin.H{
		"title": "구독관리",
	}, "subs/list.html")
}

// ShowSubsListPage 상품 리스트 페이지
func ShowSubsCancelListPage(c *gin.Context) {
	render(c, gin.H{
		"title": "구독관리(해지)",
	}, "subs/cancellist.html")
}

// ShowSubsViewPage 상품 상세 페이지
func ShowSubsViewPage(c *gin.Context) {
	subsNo, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(err)

		render(c, gin.H{
			"title": "Not Found",
		}, "error/notfound.html")

		return
	}

	info := models.GetSubsInfo(subsNo)
	if info.SubsNo < 1 {
		render(c, gin.H{
			"title": "Not Found",
		}, "error/notfound.html")
		return
	}

	typelist := models.GetCodeList("ITEM_CATEGORY", false)
	boxlist := models.GetCodeList("BOX_TYPE", false)
	weeklist := models.GetCodeList("DELIVERY_DOW", false)
	periodlist := models.GetCodeList("DELIVERY_PERIOD", false)
	addresslist := models.GetUserAddressList(info.UserNo)
	cardlist := models.GetUserCardList(info.UserNo)

	render(c, gin.H{
		"title":       "구독상세",
		"type":        typelist,
		"box":         boxlist,
		"week":        weeklist,
		"period":      periodlist,
		"addresslist": addresslist,
		"cardlist":    cardlist,
		"info":        info,
	}, "subs/view.html")
}

// GetSubsList 주문 조회
func GetSubsList(c *gin.Context) {
	userNo, _ := strconv.ParseInt(c.DefaultQuery("userno", "0"), 10, 64)
	status := strings.Split(c.DefaultQuery("status", ""), "|")

	list := models.GetSubsList(userNo, status)

	c.JSON(http.StatusOK, gin.H{"data": list})
}

// SubsManage 구독 정보 관리
type SubsManage struct {
	SubsNo      string
	CateTypes   []string
	BoxType     string
	SubsPrice   string
	PeriodDay   string
	NextDate    string
	RcvName     string
	ContactNo   string
	PostNo      string
	MainAddress string
	SubAddress  string
	ReqMsg      string
	CardRegNo   string
	AddnlDesc   string
}

// SetSubs 구독 정보 변경
func SetSubs(c *gin.Context) {
	var info SubsManage

	err := c.ShouldBind(&info)
	if err != nil {
		log.Println(err)

		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	sNo, _ := strconv.ParseInt(info.SubsNo, 10, 64)

	cType := strings.Join(info.CateTypes, "|")

	price, _ := strconv.Atoi(info.SubsPrice)
	pDay, _ := strconv.Atoi(info.PeriodDay)
	nDate, _ := time.Parse("2006-01-02", info.NextDate)

	cNo, _ := strconv.ParseInt(info.CardRegNo, 10, 64)

	result := models.SetSubs(models.SubsMst{SubsNo: sNo, CateType: cType, BoxType: info.BoxType, SubsPrice: price, PeriodDay: pDay, NextDate: &nDate,
		RcvName: &info.RcvName, ContactNo: &info.ContactNo, PostNo: &info.PostNo, MainAddress: &info.MainAddress, SubAddress: &info.SubAddress, ReqMsg: &info.ReqMsg,
		CardRegNo: cNo, AddnlDesc: &info.AddnlDesc})

	if result < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"result": result})
		return
	}

	models.AddSubsHist(sNo, "change", "구독정보 변경", GetAdminID(c))

	go common.SendSlackMessage("system", fmt.Sprintf("구독(%d) 정보변경 - %s", sNo, GetAdminID(c)))

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// CancelSubs 구독 해지
func CancelSubs(c *gin.Context) {
	subsNo, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println(err)

		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	s := models.GetSubsInfo(subsNo)

	result := models.CancelSubs(subsNo)

	if result < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"result": result})
		return
	}

	models.AddSubsHist(subsNo, "cancel", "구독 해지", GetAdminID(c))

	go common.SendSlackMessage("system", fmt.Sprintf("구독(%d) 해지 - %s", subsNo, GetAdminID(c)))

	u := models.GetUserInfo(s.UserNo)
	go common.SendAlrimTalk(u.UserPhone, "bizp_2021083018402810216575017",
		fmt.Sprintf(`%s님의 큐잇 정기구독 해지가 완료되었습니다.

신선한 농산물과 건강한 식생활이 필요할 땐, 언제든 다시 큐잇을 이용하실 수 있어요.

도움이 필요하시면 고객센터에 문의해 주세요.

고객센터 070-4166-6077

다시 만나뵐 날을 기다릴게요.
큐잇 드림`, u.UserName),
		[]common.AlrimTalkButton{})

	c.JSON(http.StatusOK, gin.H{"result": result})
}
