package models

import (
	"fmt"
	"log"
	"strings"
	"time"
	"yelloment-admin/common"
)

// OrderMst 주문
type OrderMst struct {
	OrderNo        int64      `db:"OrderNo" json:"OrderNo"`
	UserNo         int64      `db:"UserNo" json:"UserNo"`
	CardRegNo      int64      `db:"CardRegNo" json:"CardRegNo"`
	OrderType      string     `db:"OrderType" json:"OrderType"`
	CateType       string     `db:"CateType" json:"CateType"`
	BoxType        string     `db:"BoxType" json:"BoxType"`
	TagGroupNo     *int64     `db:"TagGroupNo" json:"TagGroupNo"`
	SubsNo         *int64     `db:"SubsNo" json:"SubsNo"`
	OrderRound     *int       `db:"OrderRound" json:"OrderRound"`
	OrderPrice     int        `db:"OrderPrice" json:"OrderPrice"`
	DiscountPrice  int        `db:"DiscountPrice" json:"DiscountPrice"`
	RcvName        *string    `db:"RcvName" json:"RcvName"`
	ContactNo      *string    `db:"ContactNo" json:"ContactNo"`
	MainAddress    *string    `db:"MainAddress" json:"MainAddress"`
	SubAddress     *string    `db:"SubAddress" json:"SubAddress"`
	PostNo         *string    `db:"PostNo" json:"PostNo"`
	ReqDelivDate   *time.Time `db:"ReqDelivDate" json:"ReqDelivDate"`
	ReqMsg         *string    `db:"ReqMsg" json:"ReqMsg"`
	AddnlDesc      *string    `db:"AddnlDesc" json:"AddnlDesc"`
	DelivCo        *string    `db:"DelivCo" json:"DelivCo"`
	DelivInvoiceNo *string    `db:"DelivInvoiceNo" json:"DelivInvoiceNo"`
	StatusCode     string     `db:"StatusCode" json:"StatusCode"`
	RegDate        *time.Time `db:"RegDate" json:"RegDate"`
	UpdDate        *time.Time `db:"UpdDate" json:"UpdDate"`
	CnlDate        *time.Time `db:"CnlDate" json:"CnlDate"`

	UserName   *string    `db:"UserName" json:"UserName"`
	UserEmail  *string    `db:"UserEmail" json:"UserEmail"`
	UserPhone  *string    `db:"UserPhone" json:"UserPhone"`
	UserType   *string    `db:"UserType" json:"UserType"`
	UserGender *string    `db:"UserGender" json:"UserGender"`
	BirthDay   *time.Time `db:"BirthDay" json:"BirthDay"`
	Funnel     *string    `db:"Funnel" json:"Funnel"`

	CardName   *string `db:"CardName" json:"CardName"`
	CardNumber *string `db:"CardNumber" json:"CardNumber"`
	CardKey    *string `db:"CardKey" json:"CardKey"`
	CardStatus *string `db:"CardStatus" json:"CardStatus"`

	ReviewDesc *string    `db:"ReviewDesc" json:"ReviewDesc"`
	ReviewDate *time.Time `db:"ReviewDate" json:"ReviewDate"`

	SubsSeq     int     `db:"SubsSeq" json:"SubsSeq"`
	PeriodDay   *int    `db:"PeriodDay" json:"PeriodDay"`
	ItemsString *string `db:"ItemsString" json:"ItemsString"`

	OrderTypeDesc string     `json:"OrderTypeDesc"`
	CateTypes     []string   `json:"CateTypes"`
	CateTypeDesc  string     `json:"CateTypeDesc"`
	CateInfos     []CateInfo `json:"CateInfos"`
	BoxTypeDesc   string     `json:"BoxTypeDesc"`
	StatusDesc    string     `json:"StatusDesc"`
	DelivURL      string     `json:"DelivUrl"`

	ItemList    []ItemMst
	PaymentList []OrderPayment
}

type CateInfo struct {
	CateCode string `json:"CateCode"`
	CateDesc string `json:"CateDesc"`
}

// OrderHist 주문 기록
type OrderHist struct {
	OrderNo    int64  `db:"OrderNo" json:"OrderNo"`
	StatusCode string `db:"StatusCode" json:"StatusCode"`
	HistDesc   string `db:"HistDesc" json:"HistDesc"`
	ExecUser   string `db:"ExecUser" json:"ExecUser"`
	ExecDate   string `db:"ExecDate" json:"ExecDate"`
}

func (o *OrderMst) setTypeDesc() {
	o.OrderTypeDesc = GetCodeLabel("ORDER_TYPE", o.OrderType)
	o.CateTypes = strings.Split(o.CateType, "|")
	o.CateTypeDesc = GetCodeLabelSeparate("ITEM_CATEGORY", o.CateType)
	o.BoxTypeDesc = GetCodeLabel("BOX_TYPE", o.BoxType)
	o.StatusDesc = GetCodeLabel("ORDER_STATUS", o.StatusCode)

	for _, code := range strings.Split(o.CateType, "|") {
		o.CateInfos = append(o.CateInfos, CateInfo{CateCode: code, CateDesc: GetCodeLabel("ITEM_CATEGORY", code)})
	}

	if o.DelivInvoiceNo != nil {
		switch s := *o.DelivCo; s {
		case "한진택배":
			o.DelivURL = "//www.hanjin.co.kr/kor/CMS/DeliveryMgr/WaybillResult.do?mCode=MN038&schLang=KR&wblnum=" + *o.DelivInvoiceNo
		case "롯데택배":
			o.DelivURL = "//www.lotteglogis.com/mobile/reservation/tracking/linkView?InvNo=" + *o.DelivInvoiceNo
		case "CJ대한통운":
			o.DelivURL = "//www.cjlogistics.com/ko/tool/parcel/tracking?gnbInvcNo=" + *o.DelivInvoiceNo
		case "우체국":
			o.DelivURL = "//service.epost.go.kr/trace.RetrieveDomRigiTraceList.comm?sid1=" + *o.DelivInvoiceNo
		case "로젠":
			o.DelivURL = "//www.ilogen.com/m/personal/trace/" + *o.DelivInvoiceNo
		default:
			o.DelivURL = ""
		}
	}
}

// GetOrderInfo 주문 조회
func GetOrderInfo(orderNo int64) OrderMst {
	db := common.GetDB()

	info := OrderMst{}

	err := db.Get(&info, `SELECT om.OrderNo, om.UserNo, om.CardRegNo, om.TagGroupNo, om.SubsNo
	, om.OrderType, om.CateType, om.BoxType, om.OrderRound
	, om.OrderPrice, om.DiscountPrice
	, om.RcvName, om.MainAddress, om.SubAddress, om.ContactNo, om.PostNo, om.ReqDelivDate, om.ReqMsg
	, om.AddnlDesc, om.DelivCo, om.DelivInvoiceNo, om.StatusCode, om.RegDate, om.UpdDate, om.CnlDate
	, um.UserName, um.UserType, um.UserEmail, um.UserPhone, um.UserGender, um.BirthDay, um.Funnel
	, uc.CardName, uc.CardNumber, uc.CardKey, uc.StatusCode AS CardStatus
	FROM OrderMst om
	LEFT OUTER JOIN UserMst um ON om.UserNo = um.UserNo
	LEFT OUTER JOIN UserCard uc ON om.CardRegNo = uc.CardRegNo
	WHERE OrderNo = ?`, orderNo)

	if err != nil {
		return info
	}

	info.setTypeDesc()
	info.ItemList = GetOrderItemList(info.OrderNo)
	info.PaymentList = GetOrderPaymentList(info.OrderNo)

	return info
}

// GetOrderList 주문 리스트 조회
func GetOrderList(status []string, uNo int64, sNo int64, isThisWeek bool) []OrderMst {
	db := common.GetDB()

	list := []OrderMst{}

	query := fmt.Sprintf(`WITH
	root AS (
		SELECT om.OrderNo, om.UserNo, om.CardRegNo, om.TagGroupNo, om.SubsNo
		, om.OrderType, om.CateType, om.BoxType, om.OrderRound
		, om.OrderPrice, om.DiscountPrice
		, om.RcvName, om.MainAddress, om.SubAddress, om.ContactNo, om.PostNo, om.ReqDelivDate, om.ReqMsg
		, om.AddnlDesc, om.DelivCo, om.DelivInvoiceNo, om.StatusCode, om.RegDate, om.UpdDate, om.CnlDate
		, um.UserName, um.UserType, um.UserEmail, um.UserPhone, um.UserGender, um.BirthDay, um.Funnel
		, uc.CardName, uc.CardNumber, uc.CardKey, uc.StatusCode AS CardStatus
		FROM OrderMst om
		LEFT OUTER JOIN UserMst um ON om.UserNo = um.UserNo
		LEFT OUTER JOIN UserCard uc ON om.CardRegNo = uc.CardRegNo
		WHERE 1 = 1
		-- [1] AND om.StatusCode IN ('%s')
		-- [2] AND om.UserNo = %d
		-- [3] AND om.SubsNo = %d
		-- [4] AND om.ReqDelivDate <= DATE_ADD(NOW() , INTERVAL 7 DAY)
	),
	i AS (
		SELECT root.OrderNo, GROUP_CONCAT(CONCAT(im.ItemName, ' ', im.UnitAmt, im.UnitType, ' x',oi.ItemCnt) ORDER BY im.ItemName ASC SEPARATOR '|') AS ItemsString
		FROM root
		INNER JOIN OrderItem oi ON root.OrderNo = oi.OrderNo
		INNER JOIN ItemMst im ON oi.ItemNo = im.ItemNo
		GROUP BY root.OrderNo

	),
	s AS (
		SELECT sm.SubsNo, sm.PeriodDay
		,	CASE @GROUPING WHEN sm.UserNo THEN @RANK := @RANK + 1 ELSE @RANK := 1 END AS SubsSeq
		,	@GROUPING := sm.UserNo AS GROUPING
		FROM SubsMst sm, (SELECT @GROUPING := 0, @RANK := 0) XX
		ORDER BY sm.UserNo, sm.SubsNo
	)
	SELECT root.*
	, i.ItemsString
	, IFNULL(s.SubsSeq, 0) AS SubsSeq, s.PeriodDay
	FROM root
	LEFT OUTER JOIN i on root.OrderNo = i.OrderNo
	LEFT OUTER JOIN s ON root.SubsNo = s.SubsNo
	ORDER BY root.OrderNo DESC`, strings.Join(status, "','"), uNo, sNo)

	if len(status) > 0 && len(status[0]) > 0 {
		query = strings.ReplaceAll(query, "-- [1]", "")
	}

	if uNo > 0 {
		query = strings.ReplaceAll(query, "-- [2]", "")
	}

	if sNo > 0 {
		query = strings.ReplaceAll(query, "-- [3]", "")
	}

	if isThisWeek {
		query = strings.ReplaceAll(query, "-- [4]", "")
	}

	rows, err := db.Queryx(query)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		r := OrderMst{}
		rows.StructScan(&r)

		r.setTypeDesc()

		if strings.Contains(r.StatusCode, "payment") || strings.Contains(r.StatusCode, "delivery") {
			r.PaymentList = GetOrderPaymentList(r.OrderNo)
		}

		list = append(list, r)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
	}

	return list
}

// GetOrderPrevList 과거 주문 리스트 조회
func GetOrderPrevList(orderno int64, userno int64, subsno int64) []OrderMst {
	db := common.GetDB()

	list := []OrderMst{}

	query := fmt.Sprintf(`SELECT om.OrderNo, om.UserNo, om.CardRegNo, om.TagGroupNo
	, om.SubsNo, om.OrderType, om.CateType, om.BoxType, om.OrderRound, om.OrderPrice, om.DiscountPrice, om.ReqDelivDate
	, om.RcvName, om.MainAddress, om.SubAddress, om.ContactNo, om.PostNo, om.ReqDelivDate, om.ReqMsg
	, om.DelivCo, om.DelivInvoiceNo, om.StatusCode, om.RegDate
	, rm.ReviewDesc, rm.RegDate AS ReviewDate
	FROM OrderMst om
	LEFT OUTER JOIN ReviewMst rm ON om.OrderNo = rm.OrderNo
	WHERE om.StatusCode NOT IN ('cancel')
	AND om.UserNo = %d
	AND om.OrderNo < %d
	AND om.SubsNo = %d
	ORDER BY om.OrderNo DESC`, userno, orderno, subsno)

	rows, err := db.Queryx(query)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		r := OrderMst{}
		rows.StructScan(&r)

		r.setTypeDesc()
		r.ItemList = GetOrderItemList(r.OrderNo)

		list = append(list, r)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
	}

	return list
}

// SetOrder 주문 변경
func SetOrder(s OrderMst) int64 {
	db := common.GetDB()

	if s.OrderNo < 1 {
		return -1
	}

	_, err := db.NamedExec(`UPDATE OrderMst
	SET CardRegNo=:CardRegNo, CateType=:CateType, BoxType=:BoxType
	, OrderPrice=:OrderPrice, DiscountPrice=:DiscountPrice
	, RcvName=:RcvName, MainAddress=:MainAddress, SubAddress=:SubAddress, ContactNo=:ContactNo, PostNo=:PostNo, ReqMsg=:ReqMsg
	, AddnlDesc=:AddnlDesc, UpdDate=NOW()
	WHERE OrderNo=:OrderNo;`, s)
	if err != nil {
		log.Println(err)
		return -2
	}

	return int64(s.OrderNo)
}

// CancelOrder 주문 취소
func CancelOrder(oNo int64, exUser string) int64 {
	db := common.GetDB()

	if oNo < 1 {
		return -1
	}

	_, err := db.Exec(`UPDATE OrderMst
	SET StatusCode = 'cancel'
	, CnlDate = NOW()
	WHERE OrderNo=?;`, oNo)
	if err != nil {
		log.Println(err)
		return -2
	}

	AddOrderHist(oNo, "cancel", "주문취소", exUser)

	return oNo
}

// UpdateOrderDelivInvoice 주문 송장번호 등록
func UpdateOrderDelivInvoice(oNo int64, dCo string, iNo string) int64 {
	db := common.GetDB()

	if oNo > 0 {
		res, err := db.Exec(`UPDATE OrderMst
		SET DelivCo=?, DelivInvoiceNo=?
		WHERE OrderNo=?;`, dCo, iNo, oNo)
		if err != nil {
			log.Println(err)
			return -1
		}

		ra, err := res.RowsAffected()
		if err != nil {
			log.Println(err)
			return -2
		}

		if ra != 1 {
			log.Println(ra)
			return -3
		}
	}

	return oNo
}

// SetOrderStatus 주문 상태 변경
func SetOrderStatus(oNo int64, sCode string, desc string, exUser string) int64 {
	db := common.GetDB()

	if oNo > 0 {
		res, err := db.Exec(`UPDATE OrderMst
		SET StatusCode = ?
		, UpdDate = NOW()
		WHERE OrderNo = ?;`, sCode, oNo)
		if err != nil {
			log.Println(err)
			return -1
		}

		ra, err := res.RowsAffected()
		if err != nil {
			log.Println(err)
			return -2
		}

		if ra != 1 {
			log.Println(ra)
			return -3
		}

		AddOrderHist(oNo, sCode, desc, exUser)
	}

	return oNo
}

// AddOrderHist 주문 기록 입력
func AddOrderHist(oNo int64, sCode string, hDesc string, eUser string) int64 {
	db := common.GetDB()

	_, err := db.Exec(`INSERT INTO OrderHist
	(OrderNo, StatusCode, HistDesc, ExecUSer)
	VALUES(?, ?, ?, ?);`, oNo, sCode, hDesc, eUser)
	if err != nil {
		log.Println(err)
	}

	return oNo
}
