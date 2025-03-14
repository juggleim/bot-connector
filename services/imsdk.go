package services

import (
	"bot-connector/dbs"
	"sync"

	juggleimsdk "github.com/juggleim/imserver-sdk-go"
)

var imsdkMap *sync.Map
var imLock *sync.RWMutex

func init() {
	imsdkMap = &sync.Map{}
	imLock = &sync.RWMutex{}
}

func GetImSdk(appkey string) *juggleimsdk.JuggleIMSdk {
	if val, exist := imsdkMap.Load(appkey); exist {
		return val.(*juggleimsdk.JuggleIMSdk)
	} else {
		imLock.Lock()
		defer imLock.Unlock()

		if val, exist := imsdkMap.Load(appkey); exist {
			return val.(*juggleimsdk.JuggleIMSdk)
		} else {
			dao := dbs.AppInfoDao{}
			appinfo := dao.FindByAppkey(appkey)
			if appinfo != nil {
				sdk := juggleimsdk.NewJuggleIMSdk(appkey, appinfo.AppSecret, appinfo.ApiUrl)
				imsdkMap.Store(appkey, sdk)
				return sdk
			}
			return nil
		}
	}
}
