package models

import (
	"fmt"
	"log"
	"strings"
	"time"
	"yelloment-admin/common"

	"github.com/muesli/cache2go"
)

// CodeType 코드 타입
type CodeType struct {
	CodeTypeNo   int        `db:"CodeTypeNo" json:"CodeTypeNo"`
	CodeType     string     `db:"CodeType" json:"CodeType"`
	CodeTypeDesc string     `db:"CodeTypeDesc" json:"CodeTypeDesc"`
	IsUse        bool       `db:"IsUse" json:"IsUse"`
	RegDate      *time.Time `db:"RegDate" json:"RegDate"`
	RegUser      string     `db:"RegUser" json:"RegUser"`
}

// CodeMst 코드
type CodeMst struct {
	CodeNo     int        `db:"CodeNo" json:"CodeNo"`
	CodeType   string     `db:"CodeType" json:"CodeType"`
	CodeKey    string     `db:"CodeKey" json:"CodeKey"`
	CodeLabel  string     `db:"CodeLabel" json:"CodeLabel"`
	CodeValue  *string    `db:"CodeValue" json:"CodeValue"`
	CodeValue2 *string    `db:"CodeValue2" json:"CodeValue2"`
	CodeOrder  int        `db:"CodeOrder" json:"CodeOrder"`
	IsUse      bool       `db:"IsUse" json:"IsUse"`
	RegDate    *time.Time `db:"RegDate"`
	RegUser    string     `db:"RegUser"`
	UpdDate    *time.Time `db:"UpdDate"`
	UpdUser    *string    `db:"UpdUser"`

	CodeTypeDesc *string `db:"CodeTypeDesc" json:"CodeTypeDesc"`
}

// GetCodeTypeList 코드 타입 리스트 조회
func GetCodeTypeList() []CodeType {
	db := common.GetDB()

	list := []CodeType{}

	err := db.Select(&list, `SELECT *
	FROM CodeType
	WHERE IsUse = TRUE`)
	if err != nil {
		log.Println(err)
	}

	return list
}

// GetCodeInfo 코드 조회
func GetCodeInfo(CodeNo int) CodeMst {
	db := common.GetDB()

	info := CodeMst{}

	err := db.Get(&info, `SELECT CodeNo, CodeType, CodeKey, CodeLabel, CodeValue
	, CodeValue2, CodeOrder, IsUse, RegDate, RegUser, UpdDate
	, UpdUser
	FROM CodeMst WHERE CodeNo = ?`, CodeNo)

	if err != nil {
		log.Print(err)
	}

	return info
}

// GetCodeList 코드 리스트 조회
func GetCodeList(CodeType string, isUse bool) []CodeMst {
	db := common.GetDB()

	list := []CodeMst{}

	query := fmt.Sprintf(`SELECT cm.CodeNo, cm.CodeType, cm.CodeKey, cm.CodeOrder, cm.CodeLabel, cm.CodeValue
	, cm.CodeValue2, cm.IsUse, cm.RegDate, cm.RegUser, cm.UpdDate
	, cm.UpdUser, ct.CodeTypeDesc
	FROM CodeMst cm
	INNER JOIN CodeType ct on cm.CodeType = ct.CodeType
	WHERE 1 = 1
	-- [1] AND cm.CodeType = '%s'
	-- [2] AND cm.IsUse = 1 AND ct.IsUse = 1
	ORDER BY ct.CodeTypeNo, cm.IsUse DESC, cm.CodeOrder`, CodeType)

	if len(CodeType) > 0 {
		query = strings.ReplaceAll(query, "-- [1]", "")
	}

	if isUse {
		query = strings.ReplaceAll(query, "-- [2]", "")
	}

	err := db.Select(&list, query)
	if err != nil {
		log.Println(err)
	}

	return list
}

// GetCodeLabel 코드 라벨 조회(캐시)
func GetCodeLabel(codeType string, codeKey string) string {
	cache := cache2go.Cache("myCache")

	var codelist []CodeMst

	for {
		res, err := cache.Value("codelist")

		if err != nil {
			log.Println(err)

			cache.Add("codelist", 3*time.Minute, GetCodeList("", false))
		} else {
			codelist = res.Data().([]CodeMst)
			break
		}
	}

	for _, s := range codelist {
		if s.CodeType == codeType && s.CodeKey == codeKey {
			return s.CodeLabel
		}
	}

	return codeKey
}

// GetCodeLabelSeparate 코드 라벨 조회 구분자(캐시)
func GetCodeLabelSeparate(codeType string, codekeys string) string {
	res := strings.Split(codekeys, "|")

	var str []string
	for _, s := range res {
		str = append(str, GetCodeLabel(codeType, s))
	}

	return strings.Join(str, "/")
}

// SetCode 코드 추가/수정
func SetCode(b CodeMst) int {
	db := common.GetDB()

	if b.CodeNo == 0 {
		_, err := db.NamedExec(`INSERT INTO CodeMst (CodeType, CodeKey, CodeValue, CodeValue2, CodeLabel, CodeOrder, IsUse, RegUser)
				VALUES(:CodeType, :CodeKey, :CodeValue, :CodeValue2, :CodeLabel, :CodeOrder, :IsUse, :RegUser)`, b)
		if err != nil {
			log.Println(err)
		}
	} else {
		_, err := db.NamedExec(`UPDATE CodeMst
		SET CodeType=:CodeType, CodeKey=:CodeKey, CodeValue=:CodeValue, CodeValue2=:CodeValue2, CodeLabel=:CodeLabel, CodeOrder=:CodeOrder
		, IsUse=:IsUse, UpdDate=current_timestamp(), UpdUser=:UpdUser
		WHERE CodeNo=:CodeNo`, b)
		if err != nil {
			log.Println(err)
		}
	}

	return 0
}
