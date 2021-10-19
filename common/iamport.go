package common

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

// ImpResponse 아임포트 응답
type ImpResponse struct {
	Code     int       `json:"code"`    //0이면 정상, 0아닌 값이면 message를 확인
	Message  string    `json:"message"` //code값이 0이 아닐 때, 오류 메세지를 포함
	Response ImpResult `json:"response"`
}

// ImpResult 아임포트 응답 결과
type ImpResult struct {
	AccessToken string `json:"access_token"`
	ExpiredAt   int    `json:"expired_at"`
	Now         int    `json:"now"`

	ImpUID           string `json:"imp_uid"`            //아임포트 결제 고유ID
	MerchantUID      string `json:"merchant_uid"`       //가맹점 고유ID
	CustomerUID      string `json:"customer_uid"`       //고객 고유ID(카드등록ID)
	CustomerUIDUsage string `json:"customer_uid_usage"` //고객 고유ID 사용용도

	CustomData string `json:"custom_data"` //커스텀 데이터

	Amount   int    `json:"amount"`   //결제금액
	Currency string `json:"currency"` //화폐단위

	Name       string `json:"name"`        //결제명(상품명)
	Channel    string `json:"channel"`     //결제 발생 경로
	PaidAt     int    `json:"paid_at"`     //결제일시
	PayMethod  string `json:"pay_method"`  //결제수단
	PgID       string `json:"pg_id"`       //PG사 상점ID(MID)
	PgProvider string `json:"pg_provider"` //PG사 명칭
	PgTid      string `json:"pg_tid"`      //PG사 승인ID
	ReceiptURL string `json:"receipt_url"` //영수증URL
	StartedAt  int    `json:"started_at"`  //결제시작일시
	Status     string `json:"status"`      //결제상태: ready:미결제, paid:결제완료, cancelled:결제취소, failed:결제실패
	UserAgent  string `json:"user_agent"`  //결제 디바이스 정보
	Escrow     bool   `json:"escrow"`      //에스크로결제 여부

	FailReason string `json:"fail_reason"` //실패사유
	FailedAt   int    `json:"failed_at"`   //실패일시

	CancelAmount      int         `json:"cancel_amount"`       //취소금액
	CancelHistory     []ImpCancel `json:"cancel_history"`      //취소내역
	CancelReason      string      `json:"cancel_reason"`       //취소사유
	CancelReceiptUrls []string    `json:"cancel_receipt_urls"` //취소영수증URL
	CancelledAt       int         `json:"cancelled_at"`        //취소일시

	BuyerName     string `json:"buyer_name"`     //결제자 이름
	BuyerTel      string `json:"buyer_tel"`      //결제자 전화번호
	BuyerAddr     string `json:"buyer_addr"`     //결제자 주소
	BuyerEmail    string `json:"buyer_email"`    //결제자 이메일
	BuyerPostcode string `json:"buyer_postcode"` //결제자 우편번호

	CardCode   string `json:"card_code"`   //카드사 코드
	CardType   int    `json:"card_type"`   //카드유형
	CardName   string `json:"card_name"`   //카드사 명
	CardNumber string `json:"card_number"` //카드번호
	CardQuota  int    `json:"card_quota"`  //할부개월수
	ApplyNum   string `json:"apply_num"`   //카드사 승인정보

	BankCode string `json:"bank_code"` //은행 표준코드
	BankName string `json:"bank_name"` //은행 명

	VbankCode     string `json:"vbank_code"`      //가상계좌 은행코드
	VbankDate     int    `json:"vbank_date"`      //가상계좌 마감기한
	VbankHolder   string `json:"vbank_holder"`    //가상계좌 예금주
	VbankIssuedAt int    `json:"vbank_issued_at"` //가상계좌 발급일시
	VbankName     string `json:"vbank_name"`      //가상계좌 은행명
	VbankNum      string `json:"vbank_num"`       //가상계좌 번호

	CashReceiptIssued bool `json:"cash_receipt_issued"` //현금영수증발행여부
}

// ImpCancel 아임포트 취소내역
type ImpCancel struct {
	PgTid       string `json:"vbank_num"`    // PG사 승인취소번호
	Amount      int    `json:"amount"`       // 취소 금액
	CancelledAt int    `json:"cancelled_at"` //결제취소된 시각 UNIX timestamp ,
	Reason      string `json:"reason"`       //결제취소 사유 ,
	ReceiptURL  string `json:"receipt_url"`  //취소에 대한 매출전표 확인 URL. PG사에 따라 제공되지 않는 경우도 있음
}

//아임포트 가맹점 코드
//const merchantCode string = "imp72211200"

//아임포트 API 키
const impKey string = "2280369450519081"

//아임포트 API 암호키
const impSecret string = "Sj9Zrf3HxTKSxEj6IPcFT3SCk8LvnNyg6IfyYjfvcZnJEHQmw22q6V3r3OWrBvHIdXgbgwHtv4daM3xi" //아임포트 API 암호키

// GetImpToken 토큰 가져오기
func GetImpToken() string {
	res, err := http.PostForm("https://api.iamport.kr/users/getToken", url.Values{"imp_key": {impKey}, "imp_secret": {impSecret}})
	if err != nil {
		log.Println(err)
		return ""
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return ""
	}

	var imp ImpResponse

	err = json.Unmarshal(data, &imp)

	if err != nil {
		log.Println(err)
		return ""
	}

	if imp.Code != 0 {
		log.Println(imp.Message)
		return ""
	}

	return imp.Response.AccessToken
}

// ProcImpPayment 결제 승인
func ProcImpPayment(cUID string, mUID string, amt int, title string, bName string, bEmail string, bTel string, customData string, noticeURL string) ImpResponse {
	imp := ImpResponse{}

	token := GetImpToken()
	if len(token) < 1 {
		return imp
	}

	url := "https://api.iamport.kr/subscribe/payments/again"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("customer_uid", cUID)
	_ = writer.WriteField("merchant_uid", mUID)
	_ = writer.WriteField("amount", strconv.Itoa(amt))
	_ = writer.WriteField("name", title)
	_ = writer.WriteField("buyer_name", bName)
	_ = writer.WriteField("buyer_email", bEmail)
	_ = writer.WriteField("buyer_tel", bTel)
	_ = writer.WriteField("custom_data", customData)
	_ = writer.WriteField("notice_url", noticeURL)
	err := writer.Close()
	if err != nil {
		log.Println(err)
		return imp
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err)
		return imp
	}
	req.Header.Add("Authorization", token)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return imp
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return imp
	}

	err = json.Unmarshal(body, &imp)
	if err != nil {
		log.Println(err)
		return imp
	}

	return imp
}

// ProcImpCancel 결제 취소
func ProcImpCancel(impUID string, mUID string, reason string) ImpResponse {
	imp := ImpResponse{}

	token := GetImpToken()
	if len(token) < 1 {
		return imp
	}

	url := "https://api.iamport.kr/payments/cancel"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("imp_uid", impUID)
	_ = writer.WriteField("merchant_uid", mUID)
	_ = writer.WriteField("reason", reason)
	err := writer.Close()
	if err != nil {
		log.Println(err)
		return imp
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err)
		return imp
	}
	req.Header.Add("Authorization", token)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return imp
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return imp
	}

	err = json.Unmarshal(body, &imp)
	if err != nil {
		log.Println(err)
		return imp
	}

	return imp
}
