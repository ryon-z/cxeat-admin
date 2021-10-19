package models

import (
	"fmt"
	"log"
	"os"
	"time"
	"yelloment-admin/common"
)

// OrderPayment 주문 결제
type OrderPayment struct {
	OrderNo      int64      `db:"OrderNo" json:"OrderNo"`
	PaymentType  string     `db:"PaymentType" json:"PaymentType"`
	HostData     *string    `db:"HostData" json:"HostData"`
	TID          *string    `db:"TID" json:"TID"`
	PaymentPrice int        `db:"PaymentPrice" json:"PaymentPrice"`
	StatusCode   string     `db:"StatusCode" json:"StatusCode"`
	RegDate      *time.Time `db:"RegDate" json:"RegDate"`
	CnlDate      *time.Time `db:"CnlDate" json:"CnlDate"`
	FailReason   *string    `db:"FailReason" json:"FailReason"`
}

// GetOrderPaymentList 주문 결제 정보 조회
func GetOrderPaymentList(oNo int64) []OrderPayment {
	db := common.GetDB()

	list := []OrderPayment{}

	err := db.Select(&list, `SELECT OrderNo, PaymentType, HostData, TID, PaymentPrice, StatusCode, RegDate, FailReason
		FROM OrderPayment
		WHERE OrderNo = ?`, oNo)
	if err != nil {
		log.Println(err)
	}

	return list
}

// PayOrder 주문 결제
func PayOrder(orderNo int64, exUser string) (code int, msg string) {
	info := GetOrderInfo(orderNo)

	// 등록 카드 정보 확인
	if len(*info.CardKey) < 1 || *info.CardStatus != "normal" {
		return -1, "사용할 수 없는 카드정보"
	}

	// 결제 정보 확인
	for _, v := range info.PaymentList {
		if v.StatusCode == "normal" {
			return -2, "이미 결제된 주문"
		}
	}

	// 결제 정보 설정
	prodName := ""
	if info.OrderType == "SUBS" {
		prodName = fmt.Sprintf("큐잇 구독 %d회차", *info.OrderRound)
	} else {
		prodName = "큐잇"
	}

	muid := fmt.Sprintf("ord-%d_%d", info.OrderNo, len(info.PaymentList)+1)

	//디버그 모드에서 결제ID 형식 변환
	if os.Getenv("GIN_MODE") == "debug" {
		muid = fmt.Sprintf("testord-%d_%d", info.OrderNo, len(info.PaymentList)+1)
	}

	db := common.GetDB()

	pPrice := info.OrderPrice - info.DiscountPrice

	if pPrice < 1000 {
		return -22, "최소 결제금액 미충족"
	}

	//아임포트 결제승인 시도
	res := common.ProcImpPayment(common.ChkPtr(info.CardKey).(string), muid, pPrice, prodName, common.ChkPtr(info.UserName).(string), common.ChkPtr(info.UserEmail).(string), common.ChkPtr(info.UserPhone).(string), "커스텀데이터", "")
	if res.Code != 0 || res.Response.Status != "paid" {
		// OrderPayment 입력
		_, err := db.Exec(`INSERT INTO OrderPayment (OrderNo, PaymentType, HostData, TID, PaymentPrice, StatusCode, FailReason)
			VALUES(?, ?, ?, ?, ?, 'fail', ?);`, orderNo, "CARD", "", res.Response.ImpUID, res.Response.Amount, res.Response.FailReason)
		if err != nil {
			log.Println(err)
			return -31, "결제 실패정보 입력 오류"
		}

		return -32, fmt.Sprintf("PG승인 실패 >> %v", res)
	}

	// OrderPayment 입력
	_, err := db.Exec(`INSERT INTO OrderPayment (OrderNo, PaymentType, HostData, TID, PaymentPrice, StatusCode)
		VALUES(?, ?, ?, ?, ?, 'normal');`, orderNo, "CARD", "", res.Response.ImpUID, res.Response.Amount)
	if err != nil {
		log.Println(err)
		return -4, "결제정보 입력 오류"
	}

	rs := SetOrderStatus(orderNo, "ready-delivery", "결제완료", exUser)

	if rs < 1 {
		return -5, "주문상태 변경 오류"
	}

	return 0, ""
}

// PayCancelOrder 주문 결제 취소
func PayCancelOrder(orderNo int64, exUser string) (code int, msg string) {
	db := common.GetDB()

	info := GetOrderInfo(orderNo)

	for _, op := range info.PaymentList {
		if op.StatusCode == "normal" {
			//카드결제
			if op.PaymentType == "CARD" {
				//아임포트 결제취소 시도
				res := common.ProcImpCancel(*op.TID, "", "")
				if res.Code != 0 {
					return -3, fmt.Sprintf("PG취소 실패 >> %v", res)
				}

				// OrderPayment 변경
				_, err := db.Exec(`UPDATE OrderPayment
				SET StatusCode='cancel', CnlDate=NOW()
				WHERE OrderNo = ?
				AND PaymentType = 'CARD'
				AND TID = ?`, info.OrderNo, op.TID)
				if err != nil {
					log.Println(err)
					return -4, "결제정보 변경 오류"
				}
			}
		}
	}

	rs := CancelOrder(info.OrderNo, exUser)
	if rs < 1 {
		return -5, "주문상태 변경 오류"
	}

	return 0, ""
}
