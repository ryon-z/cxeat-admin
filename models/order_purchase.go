package models

import (
	"fmt"
	"log"
	"strings"
	"yelloment-admin/common"
)

// OrderPurchase_Item 주문 발주서 (원물)
type OrderPurchase_Item struct {
	ItemNo     int64   `db:"ItemNo" json:"ItemNo"`
	ItemName   string  `db:"ItemName" json:"ItemName"`
	UnitAmt    float32 `db:"UnitAmt" json:"UnitAmt"`
	UnitType   string  `db:"UnitType" json:"UnitType"`
	SumItemCnt int     `db:"SumItemCnt" json:"SumItemCnt"`
}

// GetOrderPurchaseItem 주문 상품 조회
func GetOrderPurchaseItem(status []string) []OrderPurchase_Item {
	db := common.GetDB()

	list := []OrderPurchase_Item{}

	query := fmt.Sprintf(`SELECT oi.ItemNo, im.ItemName, im.UnitAmt, im.UnitType, SUM(oi.ItemCnt) AS SumItemCnt
	FROM OrderMst om
	INNER JOIN OrderItem oi ON om.OrderNo = oi.OrderNo
	LEFT OUTER JOIN ItemMst im ON oi.ItemNo = im.ItemNo
	WHERE 1 = 1
	-- [1] AND om.StatusCode IN ('%s')
	GROUP by oi.ItemNo, im.ItemName, im.UnitAmt, im.UnitType`, strings.Join(status, "','"))

	if len(status) > 0 && len(status[0]) > 0 {
		query = strings.ReplaceAll(query, "-- [1]", "")
	}

	err := db.Select(&list, query)
	if err != nil {
		log.Println(err)
	}

	return list
}
