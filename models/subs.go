package models

import (
	"fmt"
	"log"
	"strings"
	"time"
	"yelloment-admin/common"
)

// SubsMst 회원 구독
type SubsMst struct {
	SubsNo      int64      `db:"SubsNo" json:"SubsNo"`
	UserNo      int64      `db:"UserNo" json:"UserNo"`
	CardRegNo   int64      `db:"CardRegNo" json:"CardRegNo"`
	TagGroupNo  *int64     `db:"TagGroupNo" json:"TagGroupNo"`
	SubsType    string     `db:"SubsType" json:"SubsType"`
	CateType    string     `db:"CateType" json:"CateType"`
	BoxType     string     `db:"BoxType" json:"BoxType"`
	SubsPrice   int        `db:"SubsPrice" json:"SubsPrice"`
	PeriodDay   int        `db:"PeriodDay" json:"PeriodDay"`
	FirstDate   *time.Time `db:"FirstDate" json:"FirstDate"`
	NextDate    *time.Time `db:"NextDate" json:"NextDate"`
	AddnlDesc   *string    `db:"AddnlDesc" json:"AddnlDesc"`
	StatusCode  string     `db:"StatusCode" json:"StatusCode"`
	RegDate     *time.Time `db:"RegDate" json:"RegDate"`
	UpdDate     *time.Time `db:"UpdDate" json:"UpdDate"`
	CnlDate     *time.Time `db:"CnlDate" json:"CnlDate"`
	CnlReason   *string    `db:"CnlReason" json:"CnlReason"`
	RcvName     *string    `db:"RcvName" json:"RcvName"`
	ContactNo   *string    `db:"ContactNo" json:"ContactNo"`
	PostNo      *string    `db:"PostNo" json:"PostNo"`
	MainAddress *string    `db:"MainAddress" json:"MainAddress"`
	SubAddress  *string    `db:"SubAddress" json:"SubAddress"`
	ReqMsg      *string    `db:"ReqMsg" json:"ReqMsg"`
	SubsReason  *string    `db:"SubsReason" json:"SubsReason"`
	RefCode     *string    `db:"RefCode" json:"RefCode"`

	SubsSeq         *int       `db:"SubsSeq" json:"SubsSeq"`
	LiveSubsSeq     *int       `db:"LiveSubsSeq" json:"LiveSubsSeq"`
	MaxOrderNo      *int64     `db:"MaxOrderNo" json:"MaxOrderNo"`
	MaxOrderRound   *int       `db:"MaxOrderRound" json:"MaxOrderRound"`
	MaxReqDelivDate *time.Time `db:"MaxReqDelivDate" json:"MaxReqDelivDate"`

	User UserMst  `db:"user" json:"User"`
	Card UserCard `db:"card" json:"Card"`

	CateTypes    []string `json:"CateTypes"`
	CateTypeDesc string   `json:"CateTypeDesc"`
	BoxTypeDesc  string   `json:"BoxTypeDesc"`
	StatusDesc   string   `json:"StatusDesc"`
}

// SubsHist 구독 기록
type SubsHist struct {
	SubsNo     int64      `db:"SubsNo" json:"SubsNo"`
	StatusCode string     `db:"StatusCode" json:"StatusCode"`
	HistDesc   string     `db:"HistDesc" json:"HistDesc"`
	ExecUser   string     `db:"ExecUser" json:"ExecUser"`
	ExecDate   *time.Time `db:"ExecDate" json:"ExecDate"`
}

func (s *SubsMst) setTypeDesc() {
	s.CateTypes = strings.Split(s.CateType, "|")
	s.CateTypeDesc = GetCodeLabelSeparate("ITEM_CATEGORY", s.CateType)
	s.BoxTypeDesc = GetCodeLabel("BOX_TYPE", s.BoxType)
	s.StatusDesc = GetCodeLabel("SUBS_STATUS", s.StatusCode)
}

// GetSubsList 상품 리스트 조회
func GetSubsList(uNo int64, status []string) []SubsMst {
	db := common.GetDB()

	list := []SubsMst{}

	query := fmt.Sprintf(`WITH
	ss AS (
	SELECT sm.SubsNo
	,	CASE @GROUPING WHEN sm.UserNo THEN @RANK := @RANK + 1 ELSE @RANK := 1 END AS SubsSeq
	,	@GROUPING := sm.UserNo AS GROUPING
	FROM SubsMst sm, (SELECT @GROUPING := 0, @RANK := 0) XX
	ORDER BY sm.UserNo, sm.SubsNo
	),
	om AS (
	SELECT SubsNo, MAX(OrderNo) AS MaxOrderNo, MAX(OrderRound) AS MaxOrderRound , MAX(ReqDelivDate) AS MaxReqDelivDate
	FROM OrderMst
	WHERE OrderType = 'SUBS'
	AND StatusCode != 'cancel'
	GROUP BY SubsNo
	)
	SELECT sm.SubsNo, sm.UserNo, sm.CardRegNo, sm.TagGroupNo
	, sm.SubsType, sm.CateType, sm.BoxType, sm.SubsPrice, sm.PeriodDay
	, sm.FirstDate, sm.NextDate, sm.AddnlDesc, sm.StatusCode, sm.RegDate
	, sm.UpdDate, sm.CnlDate, sm.CnlReason, sm.SubsReason, sm.RefCode
	, sm.RcvName, sm.ContactNo, sm.MainAddress, sm.SubAddress, sm.PostNo, sm.ReqMsg
	, um.UserName 'user.UserName', um.Funnel 'user.Funnel'
	, IFNULL(ss.SubsSeq, 0) AS SubsSeq
	, om.MaxOrderNo, IFNULL(om.MaxOrderRound,0) AS MaxOrderRound, om.MaxReqDelivDate
	FROM SubsMst sm
	INNER JOIN UserMst um ON sm.UserNo = um.UserNo
	LEFT OUTER JOIN ss ON sm.SubsNo = ss.SubsNo
	LEFT OUTER JOIN om ON sm.SubsNo = om.SubsNo
	WHERE 1 = 1
	-- [1] AND sm.UserNo = %d
	-- [2] AND sm.StatusCode IN ('%s')
	ORDER BY sm.SubsNo DESC;`, uNo, strings.Join(status, "','"))

	if uNo > 0 {
		query = strings.ReplaceAll(query, "-- [1]", "")
	}

	if len(status) > 0 && len(status[0]) > 0 {
		query = strings.ReplaceAll(query, "-- [2]", "")
	}

	rows, err := db.Queryx(query)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		r := SubsMst{}
		rows.StructScan(&r)

		r.setTypeDesc()

		list = append(list, r)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
	}

	return list
}

// GetSubsInfo 구독 조회
func GetSubsInfo(subsNo int64) SubsMst {
	db := common.GetDB()

	info := SubsMst{}

	db.Get(&info, `SELECT sm.SubsNo, sm.UserNo, sm.CardRegNo, sm.TagGroupNo
	, sm.SubsType, sm.CateType, sm.BoxType, sm.SubsPrice, sm.PeriodDay
	, sm.FirstDate, sm.NextDate, sm.AddnlDesc, sm.StatusCode, sm.RegDate
	, sm.UpdDate, sm.CnlDate, sm.CnlReason, sm.SubsReason, sm.RefCode
	, sm.RcvName, sm.ContactNo, sm.MainAddress, sm.SubAddress, sm.PostNo, sm.ReqMsg
	FROM SubsMst sm
	WHERE sm.subsNo = ?`, subsNo)

	info.setTypeDesc()

	return info
}

// SetSubs 구독 변경
func SetSubs(s SubsMst) int64 {
	db := common.GetDB()

	if s.SubsNo < 1 {
		return -1
	}

	_, err := db.NamedExec(`UPDATE SubsMst
	SET CardRegNo=:CardRegNo, CateType=:CateType, BoxType=:BoxType, SubsPrice=:SubsPrice
	, PeriodDay=:PeriodDay, NextDate=:NextDate, AddnlDesc=:AddnlDesc, UpdDate=NOW()
	, RcvName=:RcvName, MainAddress=:MainAddress, SubAddress=:SubAddress, ContactNo=:ContactNo, PostNo=:PostNo, ReqMsg=:ReqMsg
	WHERE SubsNo=:SubsNo;`, s)
	if err != nil {
		log.Println(err)
	}

	return int64(s.SubsNo)
}

// CancelSubs 구독 해지
func CancelSubs(sNo int64) int64 {
	db := common.GetDB()

	if sNo < 1 {
		return -1
	}

	_, err := db.Exec(`UPDATE SubsMst
	SET StatusCode = 'cancel'
	, CnlDate = NOW()
	WHERE SubsNo=?;`, sNo)
	if err != nil {
		log.Println(err)
	}

	return sNo
}

// AddSubsHist 구독 기록 입력
func AddSubsHist(sNo int64, sCode string, hDesc string, eUser string) int64 {
	db := common.GetDB()

	_, err := db.Exec(`INSERT INTO SubsHist
	(SubsNo, StatusCode, HistDesc, ExecUSer)
	VALUES(?, ?, ?, ?);`, sNo, sCode, hDesc, eUser)
	if err != nil {
		log.Println(err)
	}

	return sNo
}

// GenSubsOrder 구독 주문 생성
func GenSubsOrder(exUser string) int64 {
	db := common.GetDB()

	list := []SubsMst{}

	query := `WITH
	sq AS (
	SELECT sm.SubsNo
	, CASE @GROUPING WHEN sm.UserNo THEN @RANK := @RANK + 1 ELSE @RANK := 1 END AS SubsSeq
	, @GROUPING := sm.UserNo AS GROUPING
	FROM SubsMst sm, (SELECT @GROUPING := 0, @RANK := 0) XX
	ORDER BY sm.UserNo, sm.SubsNo
	),
	sqq AS (
	SELECT sm.SubsNo
	, CASE @GROUPING WHEN sm.UserNo THEN @RANK := @RANK + 1 ELSE @RANK := 1 END AS LiveSubsSeq
	, @GROUPING := sm.UserNo AS GROUPING
	FROM SubsMst sm, (SELECT @GROUPING := 0, @RANK := 0) XX
	WHERE sm.StatusCode = 'normal'
	ORDER BY sm.UserNo, sm.SubsNo
	),
	om AS (
	SELECT SubsNo, MAX(OrderNo) AS MaxOrderNo, MAX(OrderRound) AS MaxOrderRound , MAX(ReqDelivDate) AS MaxReqDelivDate
	FROM OrderMst
	WHERE OrderType = 'SUBS'
	AND StatusCode != 'cancel'
	GROUP BY SubsNo
	)
	SELECT sm.SubsNo, sm.UserNo, sm.CardRegNo, sm.TagGroupNo
	, sm.SubsType, sm.CateType, sm.BoxType, sm.SubsPrice, sm.PeriodDay, sm.FirstDate
	, sm.NextDate, sm.AddnlDesc, sm.StatusCode, sm.RegDate
	, sm.RcvName, sm.ContactNo, sm.MainAddress, sm.SubAddress, sm.PostNo, sm.ReqMsg
	, IFNULL(sq.SubsSeq, 0) AS SubsSeq
	, IFNULL(sqq.LiveSubsSeq, -1) AS LiveSubsSeq
	, om.MaxOrderNo, IFNULL(om.MaxOrderRound,0) AS MaxOrderRound, om.MaxReqDelivDate
	FROM SubsMst sm
	INNER JOIN UserMst um ON sm.UserNo = um.UserNo
	LEFT OUTER JOIN sq ON sm.SubsNo = sq.SubsNo
	LEFT OUTER JOIN sqq ON sm.SubsNo = sqq.SubsNo
	LEFT OUTER JOIN om ON sm.SubsNo = om.SubsNo
	WHERE sm.StatusCode = 'normal'
	AND um.StatusCode != 'leave'
	AND sm.NextDate < DATE_ADD(NOW(), INTERVAL 6 DAY)`

	err := db.Select(&list, query)
	if err != nil {
		log.Println(err)
	}

	i := int64(0)

	for _, s := range list {

		tx := db.MustBegin()

		nxRound := 1

		if s.MaxOrderRound != nil {
			nxRound = *s.MaxOrderRound + 1
		}

		dprice := 0

		if *s.SubsSeq == 1 && nxRound == 1 {
			// 첫 구독 첫 주문 20% 할인
			dprice = s.SubsPrice / 5
		} else {
			if *s.LiveSubsSeq > 1 {
				// 복수 구독 할인(2활성구독5%, 3활성구독 10%)
				seq := *s.LiveSubsSeq

				if seq > 5 {
					seq = 5
				}

				dprice = (s.SubsPrice / 100) * (5 * (seq - 1))
			}

			if nxRound > 1 && nxRound%10 == 0 {
				// 10회차 단위 반값 할인
				dprice = dprice + (s.SubsPrice / 2)
			}
		}

		if s.SubsPrice < dprice {
			dprice = s.SubsPrice - 1000
		}

		res, err := tx.Exec(`INSERT INTO
			OrderMst (UserNo, CardRegNo, TagGroupNo, SubsNo,
				OrderType, CateType, BoxType, OrderRound, OrderPrice, DiscountPrice,
				RcvName, ContactNo, PostNo, MainAddress, SubAddress, ReqMsg,
				ReqDelivDate, AddnlDesc, StatusCode)
			VALUES(?, ?, ?, ?,
				?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?, ?,
				?, ?, ?);`,
			s.UserNo, s.CardRegNo, s.TagGroupNo, s.SubsNo,
			"SUBS", s.CateType, s.BoxType, nxRound, s.SubsPrice, dprice,
			s.RcvName, s.ContactNo, s.PostNo, s.MainAddress, s.SubAddress, s.ReqMsg,
			s.NextDate, s.AddnlDesc, "init")
		if err != nil {
			log.Println(err)
			tx.Rollback()
		}

		orderNo, err := res.LastInsertId()

		if err != nil {
			log.Println(err)
			tx.Rollback()
		}

		AddOrderHist(orderNo, "init", "주문생성", exUser)

		_, err2 := tx.Exec(`UPDATE SubsMst
		SET NextDate=?
		WHERE SubsNo=?;`, s.NextDate.AddDate(0, 0, s.PeriodDay), s.SubsNo)
		if err2 != nil {
			log.Println(err2)
			tx.Rollback()
		}

		tx.Commit()
		i++
	}

	return i
}
