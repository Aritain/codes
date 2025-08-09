package models

type Config struct {
	TGToken        string `mapstructure:"TG_TOKEN"`
	BotDebug       bool   `mapstructure:"BOT_DEBUG"`
	GCodesEndpoint string `mapstructure:"GCODES_ENDPOINT"`
	HCodesEndpoint string `mapstructure:"HCODES_ENDPOINT"`
	ZCodesEndpoint string `mapstructure:"ZCODES_ENDPOINT"`
	GCodesName     string `mapstructure:"GCODES_NAME"`
	HCodesName     string `mapstructure:"HCODES_NAME"`
	ZCodesName     string `mapstructure:"ZCODES_NAME"`
	CodesDir       string `mapstructure:"CODES_DIR"`
}

type TGMessage struct {
	TGToken string
	UserID  int64
	Text    string
}

type SavedChat struct {
	UserID    int64
	ChatPath  string
	ChatStage int8
}

// Codes structs
type CodeData struct {
	Codes []CodeBody `json:"codes"`
}

type CodeBody struct {
	Code string `json:"code"`
}

type Subscribers struct {
	Subscribers []Subscriber `json:"subscribers"`
}

type Subscriber struct {
	TGID  int64   `json:"TGID"`
	CType []CType `json:"CType"`
}

type CType struct {
	Name     string   `json:"Name"`
	TGToggle bool     `json:"TGToggle"`
	Webhooks []string `json:"Webhooks"`
}

type EndpointMap struct {
	Name string
	URL  string
}
