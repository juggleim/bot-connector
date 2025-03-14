package errs

import "fmt"

type ErrorCode int

const (
	ErrorCode_Success  ErrorCode = 0
	ErrorCode_Unknown  ErrorCode = 10000
	ErrorCode_ParamErr ErrorCode = 10001
)

type PageInfo struct {
	Page  int `json:"page"`
	Count int `json:"count"`
}
type CommonResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func (m *CommonResp) Error() string {
	return fmt.Sprintf("%d:%s", m.Code, m.Msg)
}

func GetErrorResp(code ErrorCode) *CommonResp {
	return &CommonResp{
		Code: int(code),
	}
}

func GetCommonResp(code ErrorCode, msg string, data interface{}) *CommonResp {
	return &CommonResp{
		Code: int(code),
		Msg:  msg,
		Data: data,
	}
}

func GetErrorRespWithMsg(code ErrorCode, msg string) *CommonResp {
	return &CommonResp{
		Code: int(code),
		Msg:  msg,
	}
}

func GetSuccessResp(obj interface{}) *CommonResp {
	return &CommonResp{
		Code: int(ErrorCode_Success),
		Data: obj,
	}
}
