package models

import (
	"fmt"
	"log"
	"strings"
	"time"
	"yelloment-admin/common"
)

// 회원 상담
type CounselMst struct {
	CounselNo   int64      `db:"CounselNo" json:"CounselNo"`
	UserNo      int64      `db:"UserNo" json:"UserNo"`
	SubsNo      *int64     `db:"SubsNo" json:"SubsNo"`
	OrderNo     *int64     `db:"OrderNo" json:"OrderNo"`
	CounselType *string    `db:"CounselType" json:"CounselType"`
	CounselDesc string     `db:"CounselDesc" json:"CounselDesc"`
	StatusCode  string     `db:"StatusCode" json:"StatusCode"`
	RegUser     *string    `db:"RegUser" json:"RegUser"`
	RegDate     *time.Time `db:"RegDate" json:"RegDate"`
	UpdUser     *string    `db:"UpdUser" json:"UpdUser"`
	UpdDate     *time.Time `db:"UpdDate" json:"UpdDate"`
}

// GetCounselList 상담 내역 조회
func GetCounselList(uNo int64, sNo int64, oNo int64) []CounselMst {
	db := common.GetDB()

	list := []CounselMst{}

	query := fmt.Sprintf(`SELECT CounselNo, UserNo, SubsNo, OrderNo
	, CounselType, CounselDesc, StatusCode
	, RegUser, RegDate, UpdUser, UpdDate
	FROM CounselMst
	WHERE StatusCode != 'delete'
	AND UserNo = %d
	-- [1] AND SubsNo = %d
	-- [1-2] AND SubsNo IS NULL
	-- [2] AND OrderNo = %d
	-- [2-2] AND OrderNo IS NULL
	ORDER BY CounselNo DESC`, uNo, sNo, oNo)

	if sNo > 0 {
		query = strings.ReplaceAll(query, "-- [1]", "")
	} else {
		query = strings.ReplaceAll(query, "-- [1-2]", "")
	}

	if oNo > 0 {
		query = strings.ReplaceAll(query, "-- [2]", "")
	} else {
		query = strings.ReplaceAll(query, "-- [2-2]", "")
	}

	err := db.Select(&list, query)
	if err != nil {
		log.Println(err)
	}

	return list
}

// AddCounsel 상담 추가
func AddCounsel(c CounselMst) int64 {
	db := common.GetDB()

	result, err := db.NamedExec(`INSERT INTO CounselMst (UserNo, SubsNo, OrderNo, CounselType, CounselDesc, RegUser)
								VALUES(:UserNo, :SubsNo, :OrderNo, :CounselType, :CounselDesc, :RegUser);`, c)
	if err != nil {
		log.Println(err)
	}

	no, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
	}

	return no
}

func RemoveCounsel(cNo int64, uId string) int64 {
	db := common.GetDB()

	_, err := db.Exec(`UPDATE CounselMst
		SET StatusCode='delete', UpdUser=?, UpdDate=NOW()
		WHERE CounselNo=?;
	`, uId, cNo)
	if err != nil {
		log.Println(err)
		return -1
	}

	return cNo
}
