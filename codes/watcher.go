package codes

import (
	"encoding/json"
	"fmt"
	"log"
	c "main/common"
	cfg "main/config"
	m "main/models"
	"os"
	"reflect"
	"slices"
	"time"

	d "github.com/gtuk/discordwebhook"
)

const TIMEOUT = 5

func FetchCodes() {
	config := cfg.Get().Config()
	endpoints := []m.EndpointMap{
		{Name: config.GCodesName, URL: config.GCodesEndpoint},
		{Name: config.HCodesName, URL: config.HCodesEndpoint},
		{Name: config.ZCodesName, URL: config.ZCodesEndpoint},
	}

	for {
		newCodes, name, codesResponse := ProcessCodes(config, endpoints)
		if len(newCodes) == 0 {
			time.Sleep(TIMEOUT * time.Minute)
			continue
		}

		users, err := GetCodesUsers()
		if err != nil {
			log.Println("No subscribers found, skipping")
			time.Sleep(TIMEOUT * time.Minute)
			continue
		}
		if len(newCodes) != 0 {
			message := FormatCodes(newCodes, name)
			for _, subscriber := range users.Subscribers {
				i := GetNameMap(name)
				if subscriber.CType[i].TGToggle || len(subscriber.CType[i].Webhooks) != 0 {
					go PrepMSG(subscriber, message, i)
				}

			}
			filepath := config.CodesDir + "/" + name + ".json"
			os.Remove(filepath)
			file, _ := os.Create(filepath)
			json.NewEncoder(file).Encode(codesResponse)
			file.Close()
		}
		time.Sleep(TIMEOUT * time.Minute)
	}
}

func PrepMSG(subscriber m.Subscriber, message string, index int) {
	config := cfg.Get().Config()
	if subscriber.CType[index].TGToggle {
		// TG MSG
		var tgm m.TGMessage
		tgm.TGToken = config.TGToken
		tgm.UserID = subscriber.TGID
		tgm.Text = message
		c.SendTGMessage(tgm)
	}
	for _, webhook := range subscriber.CType[index].Webhooks {
		dMessage := d.Message{
			Content: &message,
		}
		// TODO - verify that webhook still exist
		for {
			err := d.SendMessage(webhook, dMessage)
			if err == nil {
				break
			}
			log.Print(err)
			time.Sleep(TIMEOUT * time.Second)
		}
	}
}

func ProcessCodes(config m.Config, endpoints []m.EndpointMap) (newCodes []string, name string, codesResponse m.CodeData) {
	params := map[string]string{}
	headers := map[string]string{}
	var codesStored m.CodeData
	var fetchError bool
	for _, endpoint := range endpoints {
		filepath := config.CodesDir + "/" + endpoint.Name + ".json"
		codesResponse, fetchError = c.GetRequest[m.CodeData](
			endpoint.URL,
			"json",
			params, headers,
		)
		if fetchError {
			time.Sleep(TIMEOUT * time.Minute)
			continue
		}
		data, err := os.ReadFile(filepath)
		if err == nil {
			_ = json.Unmarshal(data, &codesStored)
		}
		if reflect.DeepEqual(codesResponse.Codes, codesStored.Codes) {
			time.Sleep(TIMEOUT * time.Hour)
			continue
		}
		for _, code := range codesResponse.Codes {
			if !slices.Contains(codesStored.Codes, code) {
				newCodes = append(newCodes, code.Code)
			}
		}
		if len(newCodes) != 0 {
			name = endpoint.Name
			return
		}
	}
	return
}

func FormatCodes(codes []string, name string) (codesFormatted string) {
	codesFormatted = fmt.Sprintf("New %s for %s\n\n",
		map[bool]string{true: "code", false: "codes"}[len(codes) == 1],
		name,
	)
	for _, code := range codes {
		codesFormatted += fmt.Sprintf("%s\n", code)
	}
	return
}
