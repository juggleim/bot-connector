package apimodels

type TeleBot struct {
	AppKey    string `json:"app_key"`
	TeleBotId string `json:"tele_bot_id"`
	UserId    string `json:"user_id"`
	BotToken  string `json:"bot_token"`
}
