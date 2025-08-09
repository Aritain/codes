package codes

import (
	"encoding/json"
	cfg "main/config"
	m "main/models"
	"os"
)

func GetNameMap(name string) (index int) {
	config := cfg.Get().Config()
	nameIndex := map[string]int{
		config.GCodesName: 0,
		config.HCodesName: 1,
		config.ZCodesName: 2,
	}
	return nameIndex[name]
}

func GetCodesUsers() (subscribers m.Subscribers, err error) {
	config := cfg.Get().Config()
	filepath := config.CodesDir + "/" + "subscribers.json"
	data, err := os.ReadFile(filepath)
	if err == nil {
		_ = json.Unmarshal(data, &subscribers)
	}
	return
}
