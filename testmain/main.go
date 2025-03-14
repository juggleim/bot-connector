package main

import (
	"bot-connector/services/pbobjs"
	"bot-connector/utils"
	"encoding/base64"
	"fmt"
	"time"
)

func main() {
	val := &pbobjs.ApiKey{
		Appkey:      "nsw3sue72begyv7y",
		CreatedTime: time.Now().UnixMilli(),
	}
	bs, _ := utils.PbMarshal(val)
	encodedBs, _ := utils.AesEncrypt(bs, []byte("wFkXjqmPMS43U93J"))
	wrap := &pbobjs.ApiKeyWrap{
		AppKey: "nsw3sue72begyv7y",
		Value:  encodedBs,
	}
	abs, _ := utils.PbMarshal(wrap)
	str := base64.URLEncoding.EncodeToString(abs)
	fmt.Println(str)
}
