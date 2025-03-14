package services

import (
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
					imsdk.AddBot(juggleimsdk.BotInfo{
						BotId:    senderId,
						Nickname: nickname,
						BotType:  0,
						BotConf:  fmt.Sprintf(`{"api_key":"ChBuc3czc3VlNzJiZWd5djd5GiD_F34yGv8KiCR6RBvqhVW36-u-aigZHO1KvVyqWtapEA==","webhook":"http://ec2-13-229-207-142.ap-southeast-1.compute.amazonaws.com:8070/bot-connector/telebot/events"}`),
					})
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
	fmt.Println("xx:", key)
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

type TeleBotSendBase struct {
	BotId      string `json:"bot_id"`
	ReceiverId string `json:"receiver_id"`
}

type TextMsg struct {
	TeleBotSendBase
	Text string `json:"text"`
}

func TeleBotSendText(ctx context.Context, msg *TextMsg) errs.ErrorCode {
	botId := msg.BotId
	receiverId := msg.ReceiverId
	bot := GetTeleBot(ctx, botId)
	if bot != nil {
		bot.Send(receiverId, msg.Text)
	}
	return errs.ErrorCode_Success
}
