package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"yelloment-admin/common"
	"yelloment-admin/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// ShowOrderViewPage 주문 상세 페이지
func ShowOrderViewPage(c *gin.Context) {
	orderNo, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.Error(err)

		render(c, gin.H{
			"title": "Not Found",
		}, "error/notfound.html")
		return
	}

	info := models.GetOrderInfo(orderNo)
	if info.OrderNo < 1 {
		render(c, gin.H{
			"title": "Not Found",
		}, "error/notfound.html")
		return
	}

	typelist := models.GetCodeList("ITEM_CATEGORY", false)
	boxlist := models.GetCodeList("BOX_TYPE", false)
	addresslist := models.GetUserAddressList(info.UserNo)
	cardlist := models.GetUserCardList(info.UserNo)

	render(c, gin.H{
		"title":       "주문상세",
		"type":        typelist,
		"box":         boxlist,
		"addresslist": addresslist,
		"cardlist":    cardlist,
		"info":        info,
	}, "order/view.html")
}

// ShowOrderListPage 주문 리스트 페이지
func ShowOrderListPage(c *gin.Context) {
	render(c, gin.H{
		"title": "주문관리",
	}, "order/list.html")
}

// ShowOrderItemPickupPage 주문 상품 구성 페이지
func ShowOrderItemPickupPage(c *gin.Context) {
	bundlelist := models.GetBundleList([]string{"normal"})

	render(c, gin.H{
		"title":      "주문 상품 관리",
		"bundlelist": bundlelist,
	}, "order/itempickup.html")
}

// ShowOrderItemPickupV2Page 주문 상품 구성 페이지
func ShowOrderItemPickupV2Page(c *gin.Context) {
	bundlelist := models.GetBundleList([]string{"normal"})

	render(c, gin.H{
		"title":      "주문 상품 관리",
		"bundlelist": bundlelist,
	}, "order/itempickup_v2.html")
}

// ShowOrderPaymentPage 주문 결제 관리 페이지
func ShowOrderPaymentPage(c *gin.Context) {
	render(c, gin.H{
		"title": "주문 결제 관리",
	}, "order/paymentlist.html")
}

// ShowOrderPurchaseItemPage 원물 발주서 페이지
func ShowOrderPurchaseItemPage(c *gin.Context) {
	render(c, gin.H{
		"title": "원물 발주서",
	}, "order/purchase_item.html")
}

// ShowOrderPurchaseItemUnpaidPage 원물 발주서(미결제포함) 페이지
func ShowOrderPurchaseItemUnpaidPage(c *gin.Context) {
	render(c, gin.H{
		"title": "원물 발주서(+미결제)",
	}, "order/purchase_item_unpaid.html")
}

// ShowOrderPurchaseUserPage 고객 발주서 페이지
func ShowOrderPurchaseUserPage(c *gin.Context) {
	render(c, gin.H{
		"title": "고객 발주서",
	}, "order/purchase_user.html")
}

// ShowOrderPurchaseUserUnpaidPage 고객 발주서(미결제포함) 페이지
func ShowOrderPurchaseUserUnpaidPage(c *gin.Context) {
	render(c, gin.H{
		"title": "고객 발주서(+미결제)",
	}, "order/purchase_user_unpaid.html")
}

// ShowOrderInvoicerPage 주문 운송장 관리 페이지
func ShowOrderInvoicerPage(c *gin.Context) {
	render(c, gin.H{
		"title": "운송장 관리",
	}, "order/invoicelist.html")
}

// ShowOrderDonePage 주문 완료 관리 페이지
func ShowOrderDonePage(c *gin.Context) {
	render(c, gin.H{
		"title": "배송완료 관리",
	}, "order/donelist.html")
}

// GetOrderList 주문 조회
func GetOrderList(c *gin.Context) {
	status := strings.Split(c.DefaultQuery("status", ""), "|")
	userno, _ := strconv.ParseInt(c.DefaultQuery("userno", "0"), 10, 64)
	subsno, _ := strconv.ParseInt(c.DefaultQuery("subsno", "0"), 10, 64)
	isthisweek, _ := strconv.ParseBool(c.DefaultQuery("isthisweek", "false"))

	list := models.GetOrderList(status, userno, subsno, isthisweek)

	c.JSON(http.StatusOK, gin.H{"data": list})
}

// GetOrderPrevList 지난 주문 조회
func GetOrderPrevList(c *gin.Context) {
	orderno, _ := strconv.ParseInt(c.DefaultQuery("orderno", "0"), 10, 64)
	userno, _ := strconv.ParseInt(c.DefaultQuery("userno", "0"), 10, 64)
	subsno, _ := strconv.ParseInt(c.DefaultQuery("subsno", "0"), 10, 64)
	list := models.GetOrderPrevList(orderno, userno, subsno)

	c.JSON(http.StatusOK, gin.H{"data": list})
}

// GetOrderItemList 주문 상품 조회
func GetOrderItemList(c *gin.Context) {
	orderNo, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	list := models.GetOrderItemList(orderNo)

	c.JSON(http.StatusOK, gin.H{"data": list})
}

// GetOrderPurchaseItem 주문 발주서(원물)
func GetOrderPurchaseItem(c *gin.Context) {
	status := strings.Split(c.DefaultQuery("status", ""), "|")

	list := models.GetOrderPurchaseItem(status)

	c.JSON(http.StatusOK, gin.H{"data": list})
}

// OrderManage 주문 정보 관리
type OrderManage struct {
	OrderNo       string
	CateTypes     []string
	BoxType       string
	OrderPrice    string
	DiscountPrice string

	RcvName     string
	ContactNo   string
	MainAddress string
	SubAddress  string
	PostNo      string
	ReqMsg      string

	CardRegNo string
	AddnlDesc string
}

// SetOrder 주문 정보 변경
func SetOrder(c *gin.Context) {
	adminid := sessions.Default(c).Get("adminid").(string)
	var mng OrderManage

	err := c.ShouldBind(&mng)
	if err != nil {
		log.Println(err)

		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	oNo, _ := strconv.ParseInt(mng.OrderNo, 10, 64)

	//주문조회
	old := models.GetOrderInfo(oNo)

	//취소, 배송완료, 배송중, 배송준비중 상태 변경 불가
	if old.StatusCode == "cancel" || old.StatusCode == "done" || old.StatusCode == "in-delivery" {
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "변경할 수 없는 주문 상태 입니다"})
		return
	}

	cType := strings.Join(mng.CateTypes, "|")

	price, _ := strconv.Atoi(mng.OrderPrice)
	dprice, _ := strconv.Atoi(mng.DiscountPrice)

	cNo, _ := strconv.ParseInt(mng.CardRegNo, 10, 64)

	if old.StatusCode == "ready-delivery" {
		if old.OrderPrice != price || old.DiscountPrice != dprice || old.BoxType != mng.BoxType {
			c.JSON(http.StatusOK, gin.H{"code": -2, "msg": "이미 결제된 주문의 금액정보는 변경할 수 없습니다"})
			return
		}

		if old.CardRegNo != cNo {
			c.JSON(http.StatusOK, gin.H{"code": -2, "msg": "이미 결제된 주문의 카드정보는 변경할 수 없습니다"})
			return
		}
	}

	result := models.SetOrder(models.OrderMst{
		OrderNo: oNo, CateType: cType, BoxType: mng.BoxType, OrderPrice: price, DiscountPrice: dprice,
		RcvName: &mng.RcvName, ContactNo: &mng.ContactNo,
		MainAddress: &mng.MainAddress, SubAddress: &mng.SubAddress, PostNo: &mng.PostNo,
		ReqMsg:    &mng.ReqMsg,
		CardRegNo: cNo,
		AddnlDesc: &mng.AddnlDesc})

	if result < 1 {
		c.JSON(http.StatusOK, gin.H{"code": result, "msg": "주문 변경 실패"})
		return
	}

	models.AddOrderHist(oNo, "change", "주문정보 변경", adminid)

	go common.SendSlackMessage("system", fmt.Sprintf("주문(%d) 정보변경 - %s", oNo, adminid))

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": ""})
}

// CancelOrder 주문취소
func CancelOrder(c *gin.Context) {
	orderNo, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	info := models.GetOrderInfo(orderNo)

	if info.StatusCode == "cancel" {
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "이미 취소된 주문"})
		return
	}

	if info.StatusCode == "done" || info.StatusCode == "in-delivery" {
		c.JSON(http.StatusOK, gin.H{"code": -2, "msg": "취소할 수 없는 주문 상태"})
		return
	}

	if len(info.PaymentList) > 0 {
		code, msg := models.PayCancelOrder(orderNo, GetAdminID(c))
		if code < 0 {
			go common.SendSlackMessage("system", fmt.Sprintf("주문(%d) 취소(+결제취소) 실패(%d: %s) - %s", orderNo, code, msg, GetAdminID(c)))
			c.JSON(http.StatusOK, gin.H{"code": code, "msg": msg})
			return
		}

	} else {
		result := models.CancelOrder(orderNo, GetAdminID(c))
		if result < 1 {
			go common.SendSlackMessage("system", fmt.Sprintf("주문(%d) 취소 실패 - %s", orderNo, GetAdminID(c)))
			c.JSON(http.StatusOK, gin.H{"code": result, "msg": "주문 취소 실패"})
			return
		}
	}

	go common.SendSlackMessage("system", fmt.Sprintf("주문(%d) 취소 - %s", orderNo, GetAdminID(c)))

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": ""})
}

// GenSubsOrder 구독 주문 생성
func GenSubsOrder(c *gin.Context) {

	cnt := models.GenSubsOrder(GetAdminID(c))

	go common.SendSlackMessage("system", fmt.Sprintf("구독 주문 생성(%d 건) - %s", cnt, GetAdminID(c)))

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "", "count": cnt})
}

// OrderItemManage 구성 상품 벌크
type OrderItemManage struct {
	OrderNo string       `json:"OrderNo"`
	Items   []ItemManage `json:"items"`
}

// SetOrderItem 주문 상품 등록
func SetOrderItem(c *gin.Context) {
	var info OrderItemManage

	err := c.ShouldBind(&info)
	if err != nil {
		log.Println(err)
	}

	oNo, _ := strconv.ParseInt(info.OrderNo, 0, 64)

	o := models.GetOrderInfo(oNo)

	models.InitOrderItem(oNo)

	for _, i := range info.Items {
		iNo, _ := strconv.Atoi(i.ItemNo)
		iCt, _ := strconv.Atoi(i.ItemCnt)

		models.AddOrderItem(models.OrderItem{OrderNo: oNo, ItemNo: iNo, ItemCnt: iCt})
	}

	if o.StatusCode == "init" {
		models.SetOrderStatus(oNo, "stanby-payment", "상품구성", GetAdminID(c))
	} else {
		models.SetOrderStatus(oNo, o.StatusCode, "상품구성", GetAdminID(c))
	}

	go common.SendSlackMessage("system", fmt.Sprintf("주문(%d) 상품 구성 - %s", oNo, GetAdminID(c)))

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "", "id": oNo})
}

// ProcOrderPayment 주문 결제 승인
func ProcOrderPayment(c *gin.Context) {
	orderNo, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	go common.SendSlackMessage("system", fmt.Sprintf("주문(%d) 결제 승인 시도 - %s", orderNo, GetAdminID(c)))

	o := models.GetOrderInfo(orderNo)

	code, msg := models.PayOrder(orderNo, GetAdminID(c))
	if code != 0 {
		go common.SendSlackMessage("system", fmt.Sprintf("(주문번호: %d)결제 승인 오류:\n%s\n(%d)\n - %s", orderNo, msg, code, GetAdminID(c)))
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "결제 승인 실패", "id": orderNo})
		return
	}

	SendOrderPaymentMsg(o)

	go common.SendSlackMessage("system", fmt.Sprintf("주문(%d) 결제 승인 완료 - %s", orderNo, GetAdminID(c)))

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "", "id": orderNo})
}

// ProcOrderPaymentForUnpaid 미결제 주문 결제 승인 처리
func ProcOrderPaymentForUnpaid(c *gin.Context) {

	go common.SendSlackMessage("system", fmt.Sprintf("미결제 주문 결제 승인 시도 - %s", GetAdminID(c)))

	//미결제 주문 조회
	list := models.GetOrderList([]string{"stanby-payment"}, 0, 0, true)

	i, f := 0, 0

	for _, o := range list {
		//결제 승인 처리
		code, msg := models.PayOrder(o.OrderNo, GetAdminID(c))
		if code != 0 {
			f++
			go common.SendSlackMessage("system", fmt.Sprintf("(주문번호: %d)결제 승인 오류:\n%s\n(%d)\n - %s", o.OrderNo, msg, code, GetAdminID(c)))
			continue
		}

		SendOrderPaymentMsg(o)

		i++
	}

	go common.SendSlackMessage("system", fmt.Sprintf("미결제 주문 결제 승인 %d건 완료 - %s", i, GetAdminID(c)))

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "", "ok_count": i, "fail_count": f})
}

func SendOrderPaymentMsg(o models.OrderMst) {
	u := models.GetUserInfo(o.UserNo)

	if o.OrderType == "SUBS" {
		go common.SendAlrimTalk(u.UserPhone, "bizp_2021083018285026143301022",
			fmt.Sprintf(`%s님, 안녕하세요.
식생활을 큐레이션 하는 농산물 큐레이션 정기구독 서비스, 큐잇입니다.

%d번째 큐잇 정기구독 %s(%s) 상품 구성이 완료되어, %d원이 결제 되었습니다.

상품 및 배송 정보가 맞는지 확인해 주세요.

■ 구독 상품 : %s(%s)
■ 배송 예정일 : %s
■ 배송지 : %s

배송지 변경 또는 이번 주 구독 연기를 원하신다면 고객센터로 연락주세요.

고객센터 070-4166-6077`, u.UserName, *o.OrderRound, o.BoxTypeDesc, o.CateTypeDesc, (o.OrderPrice-o.DiscountPrice), o.BoxTypeDesc, o.CateTypeDesc, o.ReqDelivDate.Format("2006-01-02"), fmt.Sprintf("%s %s", *o.MainAddress, *o.SubAddress)),
			[]common.AlrimTalkButton{})
	} else {
		//단건 주문 알림 수정 필요
		go common.SendAlrimTalk(u.UserPhone, "bizp_2021083018285026143301022",
			fmt.Sprintf(`%s님, 안녕하세요.
식생활을 큐레이션 하는 농산물 큐레이션 정기구독 서비스, 큐잇입니다.

1번째 큐잇 정기구독 %s(%s) 상품 구성이 완료되어, %d원이 결제 되었습니다.

상품 및 배송 정보가 맞는지 확인해 주세요.

■ 구독 상품 : %s(%s)
■ 배송 예정일 : %s
■ 배송지 : %s

배송지 변경 또는 이번 주 구독 연기를 원하신다면 고객센터로 연락주세요.

고객센터 070-4166-6077`, u.UserName, o.BoxTypeDesc, o.CateTypeDesc, (o.OrderPrice-o.DiscountPrice), o.BoxTypeDesc, o.CateTypeDesc, o.ReqDelivDate.Format("2006-01-02"), fmt.Sprintf("%s %s", *o.MainAddress, *o.SubAddress)),
			[]common.AlrimTalkButton{})
	}
}

// OrderInvoiceManage 구성 상품 벌크
type OrderInvoiceManage struct {
	OrderNo   string
	DelivCo   string
	InvoiceNo string
}

// UpdateOrderInvoiceNoBulk 송장번호 등록(벌크)
func UpdateOrderInvoiceNoBulk(c *gin.Context) {

	go common.SendSlackMessage("system", fmt.Sprintf("송장번호 등록(bulk) 시도 - %s", GetAdminID(c)))

	var list []OrderInvoiceManage

	err := c.ShouldBind(&list)
	if err != nil {
		log.Println(err)
	}

	ok_cnt := 0
	fail_cnt := 0

	for _, i := range list {
		orderNo, err := strconv.ParseInt(strings.TrimSpace(i.OrderNo), 10, 64)
		if err != nil {
			log.Println(err)
			fail_cnt++
			continue
		}

		delivCo := strings.TrimSpace(i.DelivCo)

		invoicerNo, err := strconv.ParseInt(strings.ReplaceAll(strings.TrimSpace(i.InvoiceNo), "-", ""), 10, 64)
		if err != nil {
			log.Println(err)
			fail_cnt++
			continue
		}

		if !(len(delivCo) > 0 && invoicerNo > 1) {
			log.Println("비정상")
			fail_cnt++
			continue
		}

		// 송장번호 등록
		r := models.UpdateOrderDelivInvoice(orderNo, delivCo, strconv.FormatInt(invoicerNo, 10))
		if r < 0 {
			fail_cnt++
			continue
		}

		// 주문 정보 조회
		o := models.GetOrderInfo(orderNo)

		if o.StatusCode == "ready-delivery" {
			models.SetOrderStatus(o.OrderNo, "in-delivery", fmt.Sprintf("송장번호 등록(%s - %s)", i.DelivCo, i.InvoiceNo), GetAdminID(c))

			if o.ReqDelivDate.After(time.Now()) {
				// 배송시작 알림 (배송예정일 이후 등록시 발송 하지 않음)
				user := models.GetUserInfo(o.UserNo)

				go common.SendAlrimTalk(user.UserPhone, "bizp_2021083018344410216148016",
					fmt.Sprintf(`%s님의 큐잇 배송이 시작되었습니다!
아래 배송 정보를 확인해 주세요.

■ 수령인 : %s
■ 택배사 : %s
■ 송장번호 : %s

%s까지 큐잇을 받지 못하시거나, 수령하신 상품에 문제가 있으면 언제든지 큐잇 고객센터로 연락주세요!
친절하게 안내해 드리겠습니다.

고객센터 070-4166-6077`, user.UserName, *o.RcvName, *o.DelivCo, i.InvoiceNo, o.ReqDelivDate.Format("2006-01-02")),
					[]common.AlrimTalkButton{{BtnName: "배송 조회", BtnType: "DS"}})
			}
		} else {
			models.AddOrderHist(o.OrderNo, o.StatusCode, fmt.Sprintf("송장번호 등록(%s - %s)", i.DelivCo, i.InvoiceNo), GetAdminID(c))
		}

		ok_cnt++
	}

	go common.SendSlackMessage("system", fmt.Sprintf("송장번호 등록 %d건 완료, %d건 실패 - %s", ok_cnt, fail_cnt, GetAdminID(c)))

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "", "ok_cnt": ok_cnt, "fail_cnt": fail_cnt})
}

// DoneOrder 단건 배송완료
func DoneOrder(c *gin.Context) {
	orderNo, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	o := models.GetOrderInfo(orderNo)

	code, msg := ProcDoneOrder(o, GetAdminID(c))

	c.JSON(http.StatusOK, gin.H{"code": code, "msg": msg})
}

// DoneOrderBulk 벌크 배송완료
func DoneOrderBulk(c *gin.Context) {
	list := models.GetOrderList([]string{"in-delivery"}, 0, 0, false)

	cnt := 0

	for _, o := range list {
		code, _ := ProcDoneOrder(o, GetAdminID(c))

		if code > 0 {
			cnt++
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "", "count": cnt})
}

// ProcDoneOrder 배송완료 처리
func ProcDoneOrder(o models.OrderMst, aid string) (code int64, msg string) {
	if o.StatusCode == "cancel" {
		return -1, "취소된 주문"
	}

	if o.StatusCode == "done" {
		return -2, "이미 배송완료 처리된 주문"
	}

	if o.StatusCode != "in-delivery" && o.StatusCode != "ready-delivery" {
		return -3, "배송완료 처리할 수 없는 상태의 주문"
	}

	r := models.SetOrderStatus(o.OrderNo, "done", "배송완료", aid)
	if r < 1 {
		return -4, "배송완료 처리 실패"
	}

	u := models.GetUserInfo(o.UserNo)

	reviewSecret := "nrX_zt3%K&2FkFy7"
	phrase := fmt.Sprintf("%d|%s", o.UserNo, reviewSecret)
	hash := sha256.Sum256([]byte(phrase))
	encoded := hex.EncodeToString(hash[:])

	// 리뷰 요청 알림 발송
	go common.SendAlrimTalk(u.UserPhone, "bizp_2021083018383326143077023",
		fmt.Sprintf(`%s님, 안녕하세요!
이번 큐잇 이용은 어떠셨나요?

받아 보신 큐잇에 대한 소감을 남겨주시면, 다음에는 더욱더 만족스러운 구성으로 찾아갈게요!

%s님의 한 마디가 큐잇이 고객님을 알아가고, 더 꼭 맞는 상품을 구성해 드리는 데에 큰 도움이 됩니다.`, u.UserName, u.UserName),
		[]common.AlrimTalkButton{{BtnName: "리뷰 작성", BtnType: "WL",
			UrlMobile: fmt.Sprintf("http://cueat.kr/web/review/insert/once?order-id=%d&cipher=%s", o.OrderNo, encoded)}})

	go common.SendSlackMessage("system", fmt.Sprintf("주문(%d) 배송완료 - %s", o.OrderNo, aid))

	return o.OrderNo, "성공"
}
