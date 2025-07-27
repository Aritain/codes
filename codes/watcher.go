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
	"strings"
	"time"
)

const TIMEOUT = 5

func FetchCodes() {
	config := cfg.Get().Config()
	endpoints := []m.EndpointMap{
		{Name: config.GCodesName, URL: config.GCodesEndpoint},
		{Name: config.HCodesName, URL: config.HCodesEndpoint},
		{Name: config.ZCodesName, URL: config.ZCodesEndpoint},
	}
	var name string

	for {

		var newCodes []string
		newCodes, name = ProcessCodes(config, endpoints)
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
			for _, user := range users.Subscribers {
				message := FormatCodes(user.UserID, newCodes, config.CodesURL)
				var tgm models.TGMessage
				tgm.TGToken = config.TGToken
				tgm.UserID = user.TGID
				tgm.Text = message
				tgm.ParseMode = "HTML"
				go common.SendTGMessage(tgm)
			}
		}
		os.Remove(filepath)
		file, _ := os.Create(filepath)
		json.NewEncoder(file).Encode(CodesResponse)
		file.Close()
		time.Sleep(TIMEOUT * time.Hour)
	}
}

func ProcessCodes(config m.Config, endpoints []m.EndpointMap) (newCodes []string, name string) {
	params := map[string]string{}
	headers := map[string]string{}
	var CodesStored m.CodeData
	var fetchError bool
	var CodesResponse m.CodeData
	for _, endpoint := range endpoints {
		filepath := config.CodesDir + "/" + endpoint.Name + ".json"
		CodesResponse, fetchError = c.GetRequest[m.CodeData](
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
			_ = json.Unmarshal(data, &CodesStored)
		}
		if reflect.DeepEqual(CodesResponse.Codes, CodesStored.Codes) {
			time.Sleep(TIMEOUT * time.Hour)
			continue
		}
		for _, code := range CodesResponse.Codes {
			if !slices.Contains(CodesStored.Codes, code) {
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

func FormatCodes(userID string, codes []string, CodesURL string) (codesFormatted string) {
	for _, code := range codes {
		fmtURL := CodesURL
		fmtURL = strings.Replace(fmtURL, "NEW_CODE", code, -1)
		fmtURL = strings.Replace(fmtURL, "USER_ID", userID, -1)
		codesFormatted += fmt.Sprintf("<a href='%s'>%s</a>\n", fmtURL, code)
	}
	codesFormatted += "\n"
	for _, code := range codes {
		codesFormatted += fmt.Sprintf("%s\n", code)
	}
	return
}
