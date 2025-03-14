package utils

import (
	"strconv"

	"google.golang.org/protobuf/proto"
)

func PbMarshal(obj proto.Message) ([]byte, error) {
	bytes, err := proto.Marshal(obj)
	return bytes, err
}

func PbUnMarshal(bytes []byte, typeScope proto.Message) error {
	err := proto.Unmarshal(bytes, typeScope)
	return err
}

func Int642String(i int64) string {
	return strconv.FormatInt(i, 10)
}

func IntPtr(v int) *int {
	return &v
}

func Int32Ptr(v int32) *int32 {
	return &v
}

func Int64Ptr(systemID int64) *int64 {
	return &systemID
}

func BoolPtr(v bool) *bool {
	return &v
}

func StringPtr(v string) *string {
	return &v
}
