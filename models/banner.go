package models

import (
	"log"
	"time"
	"yelloment-admin/common"
)

// BannerMst 배너
type BannerMst struct {
	BannerNo     int        `db:"BannerNo" json:"BannerNo"`
	BannerType   string     `db:"BannerType" json:"BannerType"`
	BannerCode   string     `db:"BannerCode" json:"BannerCode"`
	BannerTitle  string     `db:"BannerTitle" json:"BannerTitle" form:"BannerTitle"`
	BannerDesc   *string    `db:"BannerDesc" json:"BannerDesc"`
	BannerImgURL *string    `db:"BannerImgURL" json:"BannerImgURL"`
	StatusCode   string     `db:"StatusCode" json:"StatusCode"`
	RegDate      *time.Time `db:"RegDate"`
}

// GetBannerInfo 배너 조회
func GetBannerInfo(BannerNo int) BannerMst {
	db := common.GetDB()

	info := BannerMst{}

	db.Get(&info, `SELECT * FROM BannerMst WHERE BannerNo = ?`, BannerNo)

	return info
}

// GetBannerList 배너 리스트 조회
func GetBannerList() []BannerMst {
	db := common.GetDB()

	list := []BannerMst{}

	err := db.Select(&list, `SELECT * FROM BannerMst`)
	if err != nil {
		log.Println(err)
	}

	return list
}

// SetBanner 배너 추가/수정
func SetBanner(b BannerMst) int {
	db := common.GetDB()

	if b.BannerNo == 0 {
		_, err := db.NamedExec(`INSERT INTO BannerMst (BannerType, BannerCode, BannerTitle, BannerDesc, BannerImgURL, StatusCode)
				VALUES(:BannerType, :BannerCode, :BannerTitle, :BannerDesc, :BannerImgURL, :StatusCode)`, b)
		if err != nil {
			log.Println(err)
		}
	} else {
		_, err := db.NamedExec(`UPDATE BannerMst
		SET BannerType=:BannerType, BannerCode=:BannerCode, BannerTitle=:BannerTitle, BannerDesc=:BannerDesc, BannerImgURL=:BannerImgURL, StatusCode=:StatusCode
		WHERE BannerNo=:BannerNo`, b)
		if err != nil {
			log.Println(err)
		}
	}

	return 0
}
