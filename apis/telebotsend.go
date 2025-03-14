package apis

import (
	"bot-connector/errs"
	"bot-connector/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func TeleBotSendText(ctx *gin.Context) {
	var msg services.TextMsg
	if err := ctx.BindJSON(&msg); err != nil {
		ctx.JSON(http.StatusBadRequest, errs.GetErrorResp(errs.ErrorCode_ParamErr))
		return
	}
	code := services.TeleBotSendText(services.ToCtx(ctx), &msg)
	if code != errs.ErrorCode_Success {
		ctx.JSON(http.StatusOK, errs.GetErrorResp(code))
	} else {
		ctx.JSON(http.StatusOK, errs.GetSuccessResp(nil))
	}
}

func TeleTest(ctx *gin.Context) {
	var msg Event
	err := ctx.BindJSON(&msg)
	fmt.Println(err)
	bs, _ := json.Marshal(msg)
	fmt.Println(string(bs))

}

func TeleBotEvents(ctx *gin.Context) {
	var event Event
	if err := ctx.BindJSON(&event); err != nil {
		ctx.JSON(http.StatusBadRequest, errs.GetErrorResp(errs.ErrorCode_ParamErr))
		return
	}
	bs, _ := json.Marshal(event)
	fmt.Println(string(bs))
	if event.EventType == EventType_Message {
		if len(event.Payload) > 0 {
			for _, msg := range event.Payload {
				arrs := strings.Split(msg.Receiver, "_")
				if len(arrs) > 0 {
					botId := arrs[0]
					receiverId := arrs[1]
					bot := services.GetTeleBot(services.ToCtx(ctx), botId)
					if bot != nil {
						bot.Send(receiverId, msg.MsgContent)
					}
				}
			}
		}
	}
	ctx.JSON(http.StatusOK, errs.GetSuccessResp(nil))
}

type EventType string

const (
	EventType_Message EventType = "message"
)

type Event struct {
	EventType EventType   `json:"event_type"`
	Timestamp int64       `json:"timestamp"`
	Payload   []*MsgEvent `json:"payload"`
}

type MsgEvent struct {
	Sender      string       `json:"sender"`
	Receiver    string       `json:"receiver"`
	ConverType  int          `json:"conver_type"`
	MsgType     string       `json:"msg_type"`
	MsgContent  string       `json:"msg_content"`
	MsgId       string       `json:"msg_id"`
	MsgTime     int64        `json:"msg_time"`
	MentionInfo *MentionInfo `json:"mention_info"`
}
type MentionInfo struct {
	MentionType   string   `json:"mention_type"`
	TargetUserIds []string `json:"target_user_ids"`
}
