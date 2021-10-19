package models

import (
	"fmt"
	"log"
	"strings"
	"time"
	"yelloment-admin/common"
)

// ItemMst 상품
type ItemMst struct {
	ItemNo     int        `db:"ItemNo" json:"ItemNo"`
	CateCode   string     `db:"CateCode" json:"CateCode"`
	ItemName   string     `db:"ItemName" json:"ItemName"`
	DpName     *string    `db:"DpName" json:"DpName"`
	UnitType   string     `db:"UnitType" json:"UnitType"`
	UnitAmt    float32    `db:"UnitAmt" json:"UnitAmt"`
	UnitPrice  *int       `db:"UnitPrice" json:"UnitPrice"`
	StatusCode string     `db:"StatusCode" json:"StatusCode"`
	RegDate    *time.Time `db:"RegDate" json:"RegDate"`

	CateName string `db:"CateName" json:"CateName"`
	ItemCnt  *int   `db:"ItemCnt" json:"ItemCnt"`

	ReviewScore *float32 `db:"ReviewScore" json:"ReviewScore"`
}

func (i *ItemMst) setTypeDesc() {
	i.CateName = GetCodeLabel("ITEM_CATEGORY", i.CateCode)
}

// GetItemList 상품 리스트 조회
func GetItemList(status []string) []ItemMst {
	db := common.GetDB()

	list := []ItemMst{}

	query := fmt.Sprintf(`SELECT im.ItemNo, im.CateCode, im.ItemName, im.DpName, im.UnitType
	,im.UnitAmt, im.UnitPrice, im.StatusCode, im.RegDate
    FROM ItemMst im
	WHERE 1 = 1
	-- [1] AND StatusCode IN ('%s')
	ORDER BY im.StatusCode, im.ItemName, im.UnitPrice`, strings.Join(status, "','"))

	if len(status) > 0 && len(status[0]) > 0 {
		query = strings.ReplaceAll(query, "-- [1]", "")
	}

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

// GetItemInfo 상품 조회
func GetItemInfo(itemNo int) ItemMst {
	db := common.GetDB()

	info := ItemMst{}

	db.Get(&info, `SELECT im.ItemNo, im.CateCode, im.ItemName, im.DpName, im.UnitType
	, im.UnitAmt, im.UnitPrice, im.StatusCode, im.RegDate
	FROM ItemMst im WHERE im.ItemNo = ?`, itemNo)

	return info
}

// SetItem 상품 추가/수정
func SetItem(b ItemMst) int {
	db := common.GetDB()

	if b.ItemNo == 0 {
		_, err := db.NamedExec(`INSERT INTO ItemMst (CateCode, ItemName, DpName, UnitType, UnitAmt, UnitPrice, StatusCode)
                VALUES(:CateCode, :ItemName, :DpName, :UnitType, :UnitAmt, :UnitPrice, :StatusCode)`, b)
		if err != nil {
			log.Println(err)
		}
	} else {
		_, err := db.NamedExec(`UPDATE ItemMst
        SET CateCode=:CateCode, ItemName=:ItemName, DpName=:DpName, UnitType=:UnitType, UnitAmt=:UnitAmt, UnitPrice=:UnitPrice, StatusCode=:StatusCode
        WHERE ItemNo=:ItemNo`, b)
		if err != nil {
			log.Println(err)
		}
	}

	return 0
}

// BundleMst 상품 구성
type BundleMst struct {
	BundleNo   int        `db:"BundleNo" json:"BundleNo"`
	BundleType *string    `db:"BundleType" json:"BundleType"`
	BundleName string     `db:"BundleName" json:"BundleName"`
	StatusCode string     `db:"StatusCode" json:"StatusCode"`
	RegDate    *time.Time `db:"RegDate" json:"RegDate"`

	ItemList []ItemMst
}

// BundleItem 상품 구성
type BundleItem struct {
	BundleNo int `db:"BundleNo" json:"BundleNo"`
	ItemNo   int `db:"ItemNo" json:"ItemNo"`
	ItemCnt  int `db:"ItemCnt" json:"ItemCnt"`
}

// GetBundleInfo 상품 구성 조회
func GetBundleInfo(bundleNo int) BundleMst {
	db := common.GetDB()

	info := BundleMst{}

	db.Get(&info, `SELECT bm.BundleNo, bm.BundleType, bm.BundleName, bm.StatusCode, bm.RegDate
    FROM BundleMst bm WHERE bm.BundleNo = ?`, bundleNo)

	itemlist := []ItemMst{}

	query := fmt.Sprintf(`SELECT im.ItemNo, im.CateCode, im.ItemName, im.DpName, im.UnitType
	, im.UnitAmt, im.UnitPrice, im.StatusCode, im.RegDate
	, bi.ItemCnt
	FROM BundleItem bi
	INNER JOIN ItemMst im ON bi.ItemNo = im.ItemNo
	WHERE bi.BundleNo = %d`, bundleNo)

	rows, err2 := db.Queryx(query)
	if err2 != nil {
		log.Println(err2)
	}

	for rows.Next() {
		r := ItemMst{}
		rows.StructScan(&r)

		r.setTypeDesc()

		itemlist = append(itemlist, r)
	}

	info.ItemList = itemlist

	return info
}

// GetBundleList 상품 구성 조회
func GetBundleList(status []string) []BundleMst {
	db := common.GetDB()

	list := []BundleMst{}

	query := fmt.Sprintf(`SELECT bm.BundleNo, bm.BundleType, bm.BundleName, bm.StatusCode, bm.RegDate
    FROM BundleMst bm
	WHERE 1 = 1
	-- [1] AND StatusCode IN ('%s')
	ORDER BY bm.BundleNo`, strings.Join(status, "','"))

	if len(status) > 0 && len(status[0]) > 0 {
		query = strings.ReplaceAll(query, "-- [1]", "")
	}

	rows, err := db.Queryx(query)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		r := BundleMst{}
		rows.StructScan(&r)

		query2 := fmt.Sprintf(`SELECT im.ItemNo, im.CateCode, im.ItemName, im.DpName, im.UnitType
		, im.UnitAmt, im.UnitPrice, im.StatusCode, im.RegDate
		, bi.ItemCnt
	    FROM BundleItem bi
	    INNER JOIN ItemMst im ON bi.ItemNo = im.ItemNo
	    WHERE bi.BundleNo = %d`, r.BundleNo)

		rows2, err2 := db.Queryx(query2)
		if err2 != nil {
			log.Println(err2)
		}

		r.ItemList = []ItemMst{}

		for rows2.Next() {
			r2 := ItemMst{}
			rows2.StructScan(&r2)

			r2.setTypeDesc()

			r.ItemList = append(r.ItemList, r2)
		}

		list = append(list, r)
	}

	return list
}

// SetBundle 상품 구성 추가/수정
func SetBundle(b BundleMst) int64 {
	db := common.GetDB()

	if b.BundleNo == 0 {
		result, err := db.NamedExec(`INSERT INTO BundleMst (BundleName, StatusCode)
                VALUES(:BundleName, :StatusCode)`, b)
		if err != nil {
			log.Println(err)
		}

		lastID, err2 := result.LastInsertId()

		if err2 != nil {
			log.Println(err2)
		}

		return lastID
	}

	_, err := db.NamedExec(`UPDATE BundleMst
	SET BundleName=:BundleName, StatusCode=:StatusCode
	WHERE BundleNo=:BundleNo;`, b)
	if err != nil {
		log.Println(err)
	}

	return int64(b.BundleNo)
}

// RemoveBundle 번들 삭제
func RemoveBundle(bNo int) int {
	db := common.GetDB()

	if bNo < 1 {
		return -1
	}

	_, err := db.Exec(`UPDATE BundleMst
	SET StatusCode = 'remove'
	WHERE BundleNo=?;`, bNo)
	if err != nil {
		log.Println(err)
		return -2
	}

	return bNo
}

// InitBundleItem 구성 초기화
func InitBundleItem(b int) int64 {
	db := common.GetDB()

	if b > 0 {
		result := db.MustExec(`DELETE FROM BundleItem
		WHERE BundleNo=?`, b)
		if result == nil {
			log.Println("fail")
		}
	}

	return int64(b)
}

// AddBundleItem 구성 상품 추가
func AddBundleItem(b BundleItem) int64 {
	db := common.GetDB()

	if b.BundleNo > 0 {
		_, err := db.NamedExec(`INSERT INTO BundleItem (BundleNo, ItemNo, ItemCnt)
								VALUES(:BundleNo, :ItemNo, :ItemCnt);`, b)
		if err != nil {
			log.Println(err)
		}
	}

	return int64(b.BundleNo)
}
