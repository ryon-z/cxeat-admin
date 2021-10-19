package common

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type PpurioToken struct {
	AccessToken string `json:"accesstoken"`
	Type        string `json:"type"`
	Expired     string `json:"expired"`
}

type PpurioMessage struct {
	Paccount string        `json:"account"`
	PrefKey  string        `json:"refkey,omitempty"`
	Ptype    string        `json:"type"`
	Pfrom    string        `json:"from"`
	Pto      string        `json:"to"`
	Pcontent PpurioContent `json:"content"`
}

type PpurioContent struct {
	Talk AlrimTalk `json:"at"`
}

type AlrimTalk struct {
	SenderKey    string            `json:"senderkey"`
	TemplateCode string            `json:"templatecode"`
	Message      string            `json:"message"`
	Title        string            `json:"title"`
	Buttons      []AlrimTalkButton `json:"button"`
}

type AlrimTalkButton struct {
	BtnName   string `json:"name"`
	BtnType   string `json:"type"`
	UrlPC     string `json:"url_pc"`
	UrlMobile string `json:"url_mobile"`
}

const account = "pandac_dev"
const password = "panda_d2v!"
const sendNumber = "07041666077"
const senderkey = "4bd0616191f84e212cac594213350765b74ef09a"

func NewPpurio() *PpurioMessage {
	return &PpurioMessage{Paccount: account, Pfrom: sendNumber}
}

func getToken() PpurioToken {
	r := PpurioToken{}

	url := "https://api.bizppurio.com/v1/token"
	method := "POST"

	strauth := fmt.Sprintf("%s:%s", account, password)
	sEnc := b64.StdEncoding.EncodeToString([]byte(strauth))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Println(err)
		return r
	}
	req.Header.Add("Authorization", "Basic "+sEnc)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return r
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return r
	}

	json.Unmarshal(body, &r)

	return r
}

func SendAlrimTalk(phoneNo string, tmpCode string, msg string, btns []AlrimTalkButton) {
	obj := NewPpurio()

	obj.PrefKey = phoneNo
	obj.Ptype = "at"
	obj.Pto = phoneNo
	obj.Pcontent.Talk.SenderKey = senderkey
	obj.Pcontent.Talk.TemplateCode = tmpCode
	obj.Pcontent.Talk.Message = msg
	obj.Pcontent.Talk.Buttons = btns

	token := getToken()
	if len(token.AccessToken) < 1 {
		log.Println("PPURIO TOKEN INVALID")
		return
	}

	url := "https://api.bizppurio.com/v3/message"
	method := "POST"

	reqBody, err := json.Marshal(obj)

	log.Println(string(reqBody))
	if err != nil {
		log.Println(err)
		return
	}

	payload := strings.NewReader(string(reqBody))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", token.Type, token.AccessToken))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(body))
}
