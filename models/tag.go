package models

import (
	"fmt"
	"log"
	"time"
	"yelloment-admin/common"
)

type Tag struct {
	TagNo      int64      `db:"TagNo" json:"TagNo"`
	TagGroupNo int64      `db:"TagGroupNo" json:"TagGroupNo"`
	TagType    string     `db:"TagType" json:"TagType"`
	TagLabel   string     `db:"TagLabel" json:"TagLabel"`
	TagValue   *string    `db:"TagValue" json:"TagValue"`
	IsUse      bool       `db:"IsUse" json:"IsUse"`
	RegDate    *time.Time `db:"RegDate" json:"RegDate"`
	RegUser    string     `db:"RegUser" json:"RegUser"`
	UpdDate    *time.Time `db:"UpdDate" json:"UpdDate"`
	UpdUser    *string    `db:"UpdUser" json:"UpdUser"`
}

type TagGroup struct {
	TagGroupNo   int64      `db:"TagGroupNo" json:"TagGroupNo"`
	UserNo       int64      `db:"UserNo" json:"UserNo"`
	TagGroupType *string    `db:"TagGroupType" json:"TagGroupType"`
	RegDate      *time.Time `db:"RegDate" json:"RegDate"`
}

type TagInfo struct {
	TagGroupNo  int64  `db:"TagGroupNo" json:"TagGroupNo"`
	TagType     string `db:"TagType" json:"TagType"`
	TagTypeDesc string `json:"TagTypeDesc"`

	Tags []Tag `json:"Tags"`
}

func (t *TagInfo) setTypeDesc() {
	t.TagTypeDesc = GetCodeLabel("TAG_TYPE", t.TagType)
}

// GetTagList 태그 리스트 조회
func GetTagList(tgNo int64, tgType string) []Tag {
	db := common.GetDB()

	list := []Tag{}

	err := db.Select(&list, `SELECT TagNo, TagGroupNo, TagType, TagLabel, TagValue
	, IsUse, RegDate, RegUser, UpdDate, UpdUser
	FROM Tag
	WHERE TagGroupNo = ?
	AND TagType = ?`, tgNo, tgType)

	if err != nil {
		log.Println(err)
	}

	return list
}

func GetTagInfo(tgNo int64) []TagInfo {
	db := common.GetDB()

	list := []TagInfo{}

	query := fmt.Sprintf(`SELECT TagGroupNo, TagType
		FROM Tag
		WHERE TagGroupNo = %d
		GROUP BY TagGroupNo, TagType
		ORDER BY MAX(TagNo)`, tgNo)

	rows, err := db.Queryx(query)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		r := TagInfo{}
		rows.StructScan(&r)

		r.setTypeDesc()
		r.Tags = GetTagList(r.TagGroupNo, r.TagType)

		list = append(list, r)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
	}

	return list
}
