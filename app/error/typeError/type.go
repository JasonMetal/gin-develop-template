/**
* @Author: 渣渣辉
* @Description:
* @File: error
* @Module: newsPage
* @Date: 2022-11-04 09:41:03
 */
package typeError

const (
	Fail = 1000 + iota
)

var ErrorMessageList = func() map[int32]string {
	return map[int32]string{
		Fail: "无法操作",
	}
}

func GetErrorMessage(errorCode int32) string {
	if message, ok := ErrorMessageList()[errorCode]; ok {
		return message
	}

	return "请求错误"
}
