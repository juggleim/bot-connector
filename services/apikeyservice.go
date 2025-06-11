package services

import (
	"bot-connector/configures"
	"bot-connector/services/pbobjs"
	"bot-connector/utils"
	"encoding/base64"
	"time"
)

func CheckAuth(apikey string) (*pbobjs.ApiKey, error) {
	apiKeySecret := configures.Config.BotConnector.ApiKeySecret
	bs, err := base64.URLEncoding.DecodeString(apikey)
	if err != nil {
		return nil, err
	}
	decodedBs, err := utils.AesDecrypt(bs, []byte(apiKeySecret))
	if err != nil {
		return nil, err
	}
	apiKey := &pbobjs.ApiKey{}
	err = utils.PbUnMarshal(decodedBs, apiKey)
	if err != nil {
		return nil, err
	}
	return apiKey, nil
}

func GenerateApiKey(appkey, botId, userId string) (string, error) {
	val := &pbobjs.ApiKey{
		Appkey:      appkey,
		BotId:       botId,
		UserId:      userId,
		CreatedTime: time.Now().UnixMilli(),
	}
	bs, _ := utils.PbMarshal(val)
	encodedBs, err := utils.AesEncrypt(bs, []byte(configures.Config.BotConnector.ApiKeySecret))
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(encodedBs), nil
}
