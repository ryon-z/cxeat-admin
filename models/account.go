package models

import (
	"log"
	"time"
	"yelloment-admin/common"
)

// AdminAccount 관리자 계정
type AdminAccount struct {
	AdminNo        int        `db:"AdminNo" json:"AdminNo"`
	AdminID        string     `db:"AdminID" json:"AdminID"`
	AdminSecretKey string     `db:"AdminSecretKey" json:"AdminSecretKey"`
	AdminName      string     `db:"AdminName" json:"AdminName"`
	AdminGrade     int        `db:"AdminGrade" json:"AdminGrade"`
	AdminDept      string     `db:"AdminDept" json:"AdminDept"`
	RegDate        *time.Time `db:"RegDate" json:"RegDate"`
	LastLoginDate  *time.Time `db:"LastLoginDate" json:"LastLoginDate"`
	LoginFailCount int        `db:"LoginFailCount" json:"LoginFailCount"`
	StatusCode     string     `db:"StatusCode" json:"StatusCode"`
}

// GetAdminAccount 관리자 계정 조회
func GetAdminAccount(adminID string) AdminAccount {
	db := common.GetDB()

	info := AdminAccount{}

	err := db.Get(&info, `SELECT * FROM AdminAccount WHERE AdminID = ?`, adminID)
	if err != nil {
		log.Println(err)
	}

	return info
}

// IsAdminValid 관리자 계정 확인
func IsAdminValid(adminid, password string) bool {
	account := GetAdminAccount(adminid)

	if account.AdminSecretKey == password && account.StatusCode == "normal" {
		return true
	}

	return false
}
