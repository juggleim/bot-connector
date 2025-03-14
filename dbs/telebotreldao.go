package dbs

import (
	"fmt"
)

type TeleBotRelDao struct {
	ID        int64  `gorm:"primary_key"`
	AppKey    string `gorm:"app_key"`
	TeleBotId string `gorm:"tele_bot_id"`
	UserId    string `gorm:"user_id"`
	BotToken  string `gorm:"bot_token"`
}

func (rel TeleBotRelDao) TableName() string {
	return "telebotrels"
}

func (rel TeleBotRelDao) Upsert(item TeleBotRelDao) error {
	return GetDb().Exec(fmt.Sprintf("INSERT INTO %s (app_key,tele_bot_id,user_id,bot_token)VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE user_id=VALUES(user_id), bot_token=VALUES(bot_token)", rel.TableName()), item.AppKey, item.TeleBotId, item.UserId, item.BotToken).Error
}

func (rel TeleBotRelDao) QryBotsByUserId(appkey, userId string, startId, limit int64) ([]*TeleBotRelDao, error) {
	var items []*TeleBotRelDao
	err := GetDb().Where("app_key=? and user_id=? and id>?", appkey, userId, startId).Order("id asc").Limit(limit).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (rel TeleBotRelDao) FindByBotId(appkey, botId string) (*TeleBotRelDao, error) {
	var item TeleBotRelDao
	err := GetDb().Where("app_key=? and tele_bot_id=?", appkey, botId).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (rel TeleBotRelDao) QryBots(startId, limit int64) ([]*TeleBotRelDao, error) {
	var items []*TeleBotRelDao
	err := GetDb().Where("id>?", startId).Order("id asc").Limit(limit).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
