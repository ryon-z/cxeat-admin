package models

import (
	"fmt"
	"log"
	"yelloment-admin/common"
)

// OrderItem 주문 상품
type OrderItem struct {
	OrderNo int64 `db:"OrderNo" json:"OrderNo"`
	ItemNo  int   `db:"ItemNo" json:"ItemNo"`
	ItemCnt int   `db:"ItemCnt" json:"ItemCnt"`
}

// GetOrderItemList 주문 상품 조회
func GetOrderItemList(orderNo int64) []ItemMst {
	db := common.GetDB()

	list := []ItemMst{}

	query := fmt.Sprintf(`SELECT im.ItemNo, im.CateCode, im.ItemName, im.DpName, im.UnitType
	, im.UnitAmt, im.UnitPrice, im.StatusCode, im.RegDate
	, oi.ItemCnt
	, ir.ReviewScore
	FROM OrderItem oi
	LEFT OUTER JOIN ItemMst im ON oi.ItemNo = im.ItemNo
	LEFT OUTER JOIN ItemReview ir on oi.OrderNo = ir.OrderNo AND oi.ItemNo = ir.ItemNo
	WHERE oi.OrderNo = %d
	ORDER BY im.ItemName`, orderNo)

	rows, err := db.Queryx(query)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		r := ItemMst{}
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

// InitOrderItem 구성 초기화
func InitOrderItem(o int64) int64 {
	db := common.GetDB()

	if o > 0 {
		result := db.MustExec(`DELETE FROM OrderItem
		WHERE OrderNo=?`, o)
		if result == nil {
			log.Println("fail")
		}
	}

	return int64(o)
}

// AddOrderItem 구성 상품 추가
func AddOrderItem(o OrderItem) int64 {
	db := common.GetDB()

	if o.OrderNo > 0 {
		_, err := db.NamedExec(`INSERT INTO OrderItem (OrderNo, ItemNo, ItemCnt)
								VALUES(:OrderNo, :ItemNo, :ItemCnt);`, o)
		if err != nil {
			log.Println(err)
		}
	}

	return int64(o.OrderNo)
}
