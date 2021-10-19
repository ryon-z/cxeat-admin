package common

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql" //mysql driver import
	"github.com/jmoiron/sqlx"
)

// DB 객체
var DB *sqlx.DB

// DBConnect 데이터베이스 연결
func DBConnect(dbType string, dbHost string, dbName string, dbUser string, dbPassword string) *sqlx.DB {
	db, err := sqlx.Open(dbType, fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbName))

	if err != nil {
		fmt.Println("DB connection Error: ", err)
	}

	DB = db

	return DB
}

// DBDisConnect 데이터베이스 연결해제
func DBDisConnect() error {
	err := DB.Close()
	return err
}

// GetDB 데이터베이트 객체 조회
func GetDB() *sqlx.DB {
	return DB
}
