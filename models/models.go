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

type SavedChat struct {
	UserID    int64
	ChatPath  string
	ChatStage int8
}

type Webhook struct {
	UserID int64
	Type   string
	URL    string
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
	TGID int64 `json:"TelegramID"`
}

type EndpointMap struct {
	Name string
	URL  string
}
