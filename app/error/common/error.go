package common

const (
	ReqParamErr = 1001 + iota
	InitRedisErr
	EnumNotExist
	TokenInvalid
)

var ErrorMessageList = func() map[int32]string {
	return map[int32]string{
		ReqParamErr:  "参数错误",
		InitRedisErr: "初始化redis失败",
		EnumNotExist: "配置不存在",
		TokenInvalid: "请重新登录",
	}
}

func GetErrorMessage(errorCode int32) string {
	if message, ok := ErrorMessageList()[errorCode]; ok {
		return message
	}
	return "请求错误"
}
