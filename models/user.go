package models

import (
	"fmt"
	"log"
	"strings"
	"time"
	"yelloment-admin/common"
)

// UserMst 관리자 계정
type UserMst struct {
	UserNo        int64      `db:"UserNo" json:"UserNo"`
	UserType      string     `db:"UserType" json:"UserType"`
	UserID        string     `db:"UserID" json:"UserID"`
	UserSecretKey string     `db:"UserSecretKey" json:"UserSecretKey"`
	UserName      string     `db:"UserName" json:"UserName"`
	UserPhone     string     `db:"UserPhone" json:"UserPhone"`
	UserEmail     *string    `db:"UserEmail" json:"UserEmail"`
	UserGender    *string    `db:"UserGender" json:"UserGender"`
	BirthDay      *time.Time `db:"BirthDay" json:"BirthDay"`
	IsMktAgree    bool       `db:"IsMktAgree" json:"IsMktAgree"`
	Funnel        *string    `db:"Funnel" json:"Funnel"`
	StatusCode    string     `db:"StatusCode" json:"StatusCode"`
	RegDate       *time.Time `db:"RegDate" json:"RegDate"`
	LeaveDate     *time.Time `db:"LeaveDate" json:"LeaveDate"`
	LastLoginDate *time.Time `db:"LastLoginDate" json:"LastLoginDate"`
}

// GetUserInfo 유저 기본 정보 조회
func GetUserInfo(userNo int64) UserMst {
	db := common.GetDB()

	info := UserMst{}

	db.Get(&info, `SELECT UserNo, UserType, UserID, UserSecretKey, UserName
		, UserEmail, UserPhone, UserGender, BirthDay, IsMktAgree, Funnel
		, StatusCode, RegDate, LeaveDate, LastLoginDate
		FROM UserMst WHERE UserNo = ?`, userNo)

	return info
}

// GetUserList 유저 리스트 조회
func GetUserList(status []string) []UserMst {
	db := common.GetDB()

	list := []UserMst{}

	query := fmt.Sprintf(`SELECT UserNo, UserType, UserID, UserSecretKey, UserName
	, UserEmail, UserPhone, UserGender, BirthDay, IsMktAgree, Funnel
	, StatusCode, RegDate, LeaveDate, LastLoginDate
	FROM UserMst
	WHERE 1 = 1
	-- [1] AND StatusCode IN ('%s')
	ORDER BY UserNo DESC`, strings.Join(status, "','"))

	if len(status) > 0 && len(status[0]) > 0 {
		query = strings.ReplaceAll(query, "-- [1]", "")
	}

	err := db.Select(&list, query)
	if err != nil {
		log.Println(err)
	}

	return list
}

// UserAddress 회원 주소
type UserAddress struct {
	AddressNo    int64      `db:"AddressNo" json:"AddressNo"`
	UserNo       int64      `db:"UserNo" json:"UserNo"`
	AddressLabel *string    `db:"AddressLabel" json:"AddressLabel"`
	RcvName      *string    `db:"RcvName" json:"RcvName"`
	RoadAddress  *string    `db:"RoadAddress" json:"RoadAddress"`
	LotAddress   *string    `db:"LotAddress" json:"LotAddress"`
	SubAddress   *string    `db:"SubAddress" json:"SubAddress"`
	PostNo       *string    `db:"PostNo" json:"PostNo"`
	ContactNo    *string    `db:"ContactNo" json:"ContactNo"`
	ReqMsg       *string    `db:"ReqMsg" json:"ReqMsg"`
	IsBasic      bool       `db:"IsBasic" json:"IsBasic"`
	StatusCode   string     `db:"StatusCode" json:"StatusCode"`
	RegDate      *time.Time `db:"RegDate" json:"RegDate"`
}

func GetUserAddress(addressNo int64) UserAddress {
	db := common.GetDB()

	info := UserAddress{}

	db.Get(&info, `SELECT ua.AddressNo, ua.UserNo, ua.AddressLabel, ua.RcvName, ua.RoadAddress
	, ua.LotAddress, ua.SubAddress, ua.PostNo, ua.ContactNo, ua.ReqMsg
	, ua.IsBasic, ua.StatusCode, ua.RegDate
	FROM UserAddress ua
	WHERE AddressNo = ?`, addressNo)

	return info
}

// GetUserAddressList 회원 주소 리스트 조회
func GetUserAddressList(userNo int64) []UserAddress {
	db := common.GetDB()

	list := []UserAddress{}

	err := db.Select(&list, `SELECT ua.AddressNo, ua.UserNo, ua.AddressLabel, ua.RcvName, ua.RoadAddress
		, ua.LotAddress, ua.SubAddress, ua.PostNo, ua.ContactNo, ua.ReqMsg
		, ua.IsBasic, ua.StatusCode, ua.RegDate
		FROM UserAddress ua
		WHERE UserNo = ?
		AND StatusCode = 'normal'
		ORDER BY ua.IsBasic DESC, ua.AddressNo ASC`, userNo)

	if err != nil {
		log.Println(err)
	}

	return list
}

// UserCard 관리자 계정
type UserCard struct {
	CardRegNo    int64      `db:"CardRegNo" json:"CardRegNo"`
	UserNo       int64      `db:"UserNo" json:"UserNo"`
	CardName     *string    `db:"CardName" json:"CardName"`
	CardNickName *string    `db:"CardNickName" json:"CardNickName"`
	CardNumber   *string    `db:"CardNumber" json:"CardNumber"`
	CardKey      string     `db:"CardKey" json:"CardKey"`
	IsBasic      bool       `db:"IsBasic" json:"IsBasic"`
	StatusCode   string     `db:"StatusCode" json:"StatusCode"`
	RegDate      *time.Time `db:"RegDate" json:"RegDate"`
}

// GetUserCardList 회원 주소 리스트 조회
func GetUserCardList(userNo int64) []UserCard {
	db := common.GetDB()

	list := []UserCard{}

	err := db.Select(&list, `SELECT CardRegNo, UserNo, CardName, CardNickName, CardNumber
		, CardKey, IsBasic, StatusCode, RegDate
		FROM UserCard
		WHERE UserNo = ?
		ORDER BY IsBasic DESC, StatusCode DESC, CardRegNo DESC `, userNo)

	if err != nil {
		log.Println(err)
	}

	return list
}
