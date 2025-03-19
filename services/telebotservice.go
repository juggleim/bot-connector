package services

import (
	"bot-connector/apimodels"
	"bot-connector/configures"
	"bot-connector/dbs"
	"bot-connector/errs"
	"bot-connector/utils"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	juggleimsdk "github.com/juggleim/imserver-sdk-go"
	tele "gopkg.in/telebot.v4"
)

func InitTeleBots() {
	dao := dbs.TeleBotRelDao{}
	var start int64 = 0
	for {
		bots, err := dao.QryBots(start, 100)
		if err != nil {
			break
		}
		for _, bot := range bots {
			if start < bot.ID {
				start = bot.ID
			}
			ctx := context.Background()
			ctx = context.WithValue(ctx, CtxKey_AppKey, bot.AppKey)
			GetTeleBot(ctx, bot.TeleBotId)
		}
		if len(bots) < 100 {
			break
		}
	}
}

type TeleBot struct {
	AppKey      string
	BotId       string
	UserId      string
	botInstance *tele.Bot
	isStarted   bool
	senderIdMap *sync.Map
}

func (bot *TeleBot) Start() {
	if !bot.isStarted {
		bot.isStarted = true
		bot.botInstance.Handle(tele.OnText, func(ctx tele.Context) error {
			txt := ctx.Text()
			sender := ctx.Sender()
			senderId := bot.BotId + "_" + utils.Int642String(sender.ID)
			imsdk := GetImSdk(bot.AppKey)
			if imsdk != nil {
				//init telegram userinfo
				if _, exist := bot.senderIdMap.LoadOrStore(senderId, true); !exist {
					nickname := sender.Username
					if sender.FirstName != "" {
						nickname = sender.FirstName
						if sender.LastName != "" {
							nickname = nickname + " " + sender.LastName
						}
					}
					apiKey, err := GenerateApiKey(bot.AppKey, bot.BotId, utils.Int642String(sender.ID))
					if err == nil {
						imsdk.AddBot(juggleimsdk.BotInfo{
							BotId:    senderId,
							Nickname: nickname,
							BotType:  utils.IntPtr(0),
							BotConf:  fmt.Sprintf(`{"api_key":"%s","webhook":"%s/bot-connector/telebot/events"}`, configures.Config.Domain, apiKey),
							ExtFields: map[string]string{
								"user_tag": "telegram_bot",
							},
						})
					}
				}
				imsdk.SendPrivateMsg(juggleimsdk.Message{
					SenderId:       senderId,
					TargetIds:      []string{bot.UserId},
					MsgType:        "jg:text",
					MsgContent:     fmt.Sprintf(`{"content":"%s"}`, txt),
					IsNotifySender: utils.BoolPtr(false),
				})
			}
			return nil
		})
		bot.botInstance.Start()
	}
}

func (bot *TeleBot) Stop() {
	if bot.isStarted {
		bot.botInstance.Stop()
	}
}

type TeleUser struct {
	UserId string
}

func (u *TeleUser) Recipient() string {
	return u.UserId
}

func (bot *TeleBot) Send(receiverId string, msg interface{}) (*tele.Message, error) {
	if bot.isStarted {
		return bot.botInstance.Send(&TeleUser{UserId: receiverId}, msg)
	}
	return nil, fmt.Errorf("bot[%s] not started.", bot.BotId)
}

var botCache *sync.Map

func init() {
	botCache = &sync.Map{}
}

func GetTeleBot(ctx context.Context, botId string) *TeleBot {
	appkey := GetAppKeyFromCtx(ctx)
	key := strings.Join([]string{appkey, botId}, "_")
	cacheBot, exist := botCache.LoadOrStore(key, &TeleBot{AppKey: appkey, BotId: botId, senderIdMap: &sync.Map{}})
	if exist {
		fmt.Println(cacheBot)
		return cacheBot.(*TeleBot)
	} else {
		//init
		ret := cacheBot.(*TeleBot)
		teleBotDao := dbs.TeleBotRelDao{}
		rel, err := teleBotDao.FindByBotId(appkey, botId)
		if err == nil && rel != nil {
			ret.UserId = rel.UserId
			pref := tele.Settings{
				Token: rel.BotToken,
				Poller: &tele.LongPoller{
					Timeout: 5 * time.Second,
				},
			}
			bot, err := tele.NewBot(pref)
			if err == nil {
				ret.botInstance = bot
				go func() {
					ret.Start()
				}()
			} else {
				fmt.Println("err:", err)
			}
		}
		return ret
	}
}

func RemoveTeleBot(ctx context.Context, botId string) {
	appkey := GetAppKeyFromCtx(ctx)
	key := strings.Join([]string{appkey, botId}, "_")
	if val, exist := botCache.LoadAndDelete(key); exist {
		bot := val.(*TeleBot)
		bot.Stop()
	}
}

func TeleBotAdd(ctx context.Context, req *apimodels.TeleBot) errs.ErrorCode {
	start := time.Now()
	dao := dbs.TeleBotRelDao{}
	err := dao.Upsert(dbs.TeleBotRelDao{
		AppKey:    req.AppKey,
		TeleBotId: req.TeleBotId,
		BotToken:  req.BotToken,
		UserId:    req.UserId,
	})
	fmt.Println("after db:", time.Since(start))
	if err == nil {
		RemoveTeleBot(ctx, req.TeleBotId)
		fmt.Println("after remove cache:", time.Since(start))
		GetTeleBot(ctx, req.TeleBotId)
		fmt.Println("after add cache:", time.Since(start))
	}
	return errs.ErrorCode_Success
}

func TeleBotDel(ctx context.Context, req *apimodels.TeleBot) errs.ErrorCode {
	dao := dbs.TeleBotRelDao{}
	err := dao.Delete(req.AppKey, req.TeleBotId)
	if err == nil {
		RemoveTeleBot(ctx, req.TeleBotId)
	}
	return errs.ErrorCode_Success
}
