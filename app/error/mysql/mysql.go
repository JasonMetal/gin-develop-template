package mysql

const (
	QueryErr = 100
)

var ErrorMessageList = func() map[int32]string {
	return map[int32]string{
		QueryErr: "查询错误",
	}
}

func GetErrorMessage(errorCode int32) string {
	if message, ok := ErrorMessageList()[errorCode]; ok {
		return message
	}

	return "请求错误"
}
