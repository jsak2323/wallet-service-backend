package telegram

import (
	"encoding/json"
	
	"gopkg.in/resty.v0"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
)

type SendMessageRes struct {

}

func SendMessage(message string) (RES SendMessageRes, err error) {
	token := config.CONF.TelegramBotToken

	res, err := resty.R().
		SetQueryParams(map[string]string{
			"chat_id": config.CONF.TelegramAlertChatId,
			"text": message,
		}).
		Get("https://api.telegram.org/bot"+token+"/sendMessage")
	
	if err != nil {
        return SendMessageRes{}, err
    }

	if err = json.Unmarshal(res.Body, &RES); err != nil {
		return SendMessageRes{}, err
	}
	
	return RES, nil
}