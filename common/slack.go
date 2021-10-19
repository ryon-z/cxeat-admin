package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func SendSlackMessage(channel string, msg string) string {
	url := ""

	if os.Getenv("GIN_MODE") == "debug" {
		url = "https://hooks.slack.com/services/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"
	} else {
		switch channel {
		case "system":
			url = "https://hooks.slack.com/services/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"
		case "operation":
			url = "https://hooks.slack.com/services/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"
		default:
			url = "https://hooks.slack.com/services/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"
		}
	}

	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{"text":"[Admin] %s"}`, msg))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err)
		return err.Error()
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return string(body)
}
