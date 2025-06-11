package dbs

import "time"

type AppInfoDao struct {
	ID          int64     `gorm:"primary_key"`
	AppName     string    `gorm:"app_name"`
	AppKey      string    `gorm:"app_key"`
	AppSecret   string    `gorm:"app_secret"`
	AppStatus   int       `gorm:"app_status"`
	CreatedTime time.Time `gorm:"created_time"`
	UpdatedTime time.Time `gorm:"updated_time"`
}

func (app AppInfoDao) TableName() string {
	return "apps"
}

func (app AppInfoDao) Create(item AppInfoDao) error {
	err := GetDb().Create(&item).Error
	return err
}

func (app AppInfoDao) FindByAppkey(appkey string) *AppInfoDao {
	var appItem AppInfoDao
	err := GetDb().Where("app_key=?", appkey).Take(&appItem).Error
	if err != nil {
		return nil
	}
	return &appItem
}
