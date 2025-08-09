package common

import (
	"encoding/json"
	"encoding/xml"
	"log"
	m "main/models"
	"net/http"
	"time"

	t "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const TIMEOUT = 5

func CompileYesNoKeyboard() t.ReplyKeyboardMarkup {
	var keyboard = t.NewReplyKeyboard(
		t.NewKeyboardButtonRow(
			t.NewKeyboardButton("Cancel"),
		),
	)
	return keyboard
}

func GetRequest[T any](url string, mode string, params map[string]string, headers map[string]string) (T, bool) {
	var results T

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return results, true
	}

	reqParams := req.URL.Query()
	for k, v := range params {
		reqParams.Add(k, v)
	}
	req.URL.RawQuery = reqParams.Encode()

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return results, true
	}

	if mode == "xml" {
		err = xml.NewDecoder(resp.Body).Decode(&results)
	} else if mode == "json" {
		err = json.NewDecoder(resp.Body).Decode(&results)
	}
	if err != nil {
		return results, true
	}

	return results, false
}

func SendTGMessage(tgm m.TGMessage) {
	bot, _ := t.NewBotAPI(tgm.TGToken)
	msg := t.NewMessage(tgm.UserID, tgm.Text)
	var err error

	for {
		_, err = bot.Send(msg)
		if err == nil {
			break
		}
		log.Print(err)
		time.Sleep(TIMEOUT * time.Second)
	}
}
