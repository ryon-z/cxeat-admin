package models

import (
	"fmt"
	"log"
	"strings"
	"time"
	"yelloment-admin/common"
)

// OrderItem 주문 상품
type ReviewMst struct {
	ReviewNo   int64      `db:"ReviewNo" json:"ReviewNo"`
	OrderNo    int64      `db:"OrderNo" json:"OrderNo"`
	ReviewDesc string     `db:"ReviewDesc" json:"ReviewDesc"`
	RegDate    *time.Time `db:"RegDate" json:"RegDate"`
	UpdDate    *time.Time `db:"UpdDate" json:"UpdDate"`

	SubsNo *int64 `db:"SubsNo" json:"SubsNo"`
	UserNo int64  `db:"UserNo" json:"UserNo"`

	ItemReview []ItemReview `db:"ItemReview" json:"ItemReview"`
}

// ItemReview
type ItemReview struct {
	ItemReviewNo int64      `db:"ItemReviewNo" json:"ItemReviewNo"`
	ReviewNo     int64      `db:"ReviewNo" json:"ReviewNo"`
	OrderNo      int64      `db:"OrderNo" json:"OrderNo"`
	ItemNo       int        `db:"ItemNo" json:"ItemNo"`
	ReviewScore  float32    `db:"ReviewScore" json:"ReviewScore"`
	ReviewDesc   string     `db:"ReviewDesc" json:"ReviewDesc"`
	RegDate      *time.Time `db:"RegDate" json:"RegDate"`
	UpdDate      *time.Time `db:"UpdDate" json:"UpdDate"`

	ItemName string  `db:"ItemName" json:"ItemName"`
	DpName   string  `db:"DpName" json:"DpName"`
	UnitAmt  float32 `db:"UnitAmt" json:"UnitAmt"`
	UnitType string  `db:"UnitType" json:"UnitType"`
}

func GetReviewItemList(rNo int64) []ItemReview {
	db := common.GetDB()

	list := []ItemReview{}

	err := db.Select(&list, `SELECT ir.ItemReviewNo, ir.ReviewNo, ir.OrderNo, ir.ItemNo, ir.ReviewScore, ir.ReviewDesc, ir.RegDate, ir.UpdDate
		, im.ItemName, im.DpName, im.UnitAmt, im.UnitType
		FROM ItemReview ir
		INNER JOIN ItemMst im ON ir.ItemNo = im.ItemNo
		WHERE ir.ReviewNo = ?`, rNo)
	if err != nil {
		log.Println(err)
	}

	return list
}

func GetReviewList(uNo int64, sNo int64, oNo int64) []ReviewMst {
	db := common.GetDB()

	list := []ReviewMst{}

	query := fmt.Sprintf(`SELECT rm.ReviewNo, rm.OrderNo, rm.ReviewDesc, rm.RegDate, rm.UpdDate
	, om.SubsNo, om.UserNo
	FROM ReviewMst rm
	LEFT OUTER JOIN OrderMst om ON rm.OrderNo = om.OrderNo
	WHERE 1 = 1
	-- [1] AND om.UserNo = %d
	-- [2] AND om.SubsNo = %d
	-- [3] AND rm.OrderNo = %d
	ORDER BY rm.OrderNo DESC, rm.RegDate DESC`, uNo, sNo, oNo)

	if uNo > 0 {
		query = strings.ReplaceAll(query, "-- [1]", "")
	}

	if sNo > 0 {
		query = strings.ReplaceAll(query, "-- [2]", "")
	}

	if oNo > 0 {
		query = strings.ReplaceAll(query, "-- [3]", "")
	}

	rows, err := db.Queryx(query)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		r := ReviewMst{}
		rows.StructScan(&r)

		r.ItemReview = GetReviewItemList(r.ReviewNo)

		list = append(list, r)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
	}

	return list
}
