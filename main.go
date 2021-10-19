package main

import (
	"log"
	"net/http"
	"os"
	"yelloment-admin/common"
	"yelloment-admin/router"
	"yelloment-admin/views"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// 환경변수 불러오기
func loadEnv() {
	var err error

	if len(os.Args) == 2 {
		err = godotenv.Load(os.Args[1])
	} else {
		err = godotenv.Load()
	}

	if err != nil {
		log.Println(".env 파일을 불러오지 못했습니다 : " + err.Error())
	}
}

func init() {
	// 타임존 설정
	os.Setenv("TZ", "Asia/Seoul")

	// 환경변수 불러오기
	loadEnv()

	// DB 커넥션 설정
	common.DBConnect(os.Getenv("DB_TYPE"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
}

func main() {
	gin.SetMode(os.Getenv("GIN_MODE"))

	r := gin.Default()

	//쿠키 설정
	store := cookie.NewStore([]byte("secret"))

	// 세션쿠키 설정
	r.Use(sessions.Sessions("ymsession", store))

	// 최대 전송량 설정
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	// 정적 파일 설정
	r.StaticFS("/contents", http.Dir("contents"))

	// 템플릿 함수 설정
	r.SetFuncMap(views.FuncMap())

	// 템블릿 파일 설정
	r.LoadHTMLGlob("views/**/*")

	// 라우트 설정
	router.InitRoutes(r)

	r.Run(":" + os.Getenv("PORT"))
}
