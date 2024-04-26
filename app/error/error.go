package error

import (
	"fmt"
	"strings"
)

const (
	ErrService ErrComponent = "service"
	ErrLogic   ErrComponent = "logic"
	ErrModel   ErrComponent = "model"
)

type ErrComponent string

type Error interface {
	error
	Code() int32
	SetCodeMsg(code int32)
	Message() string
	SetMessage(msg string)
	Cause() error
	SetCause(err error)
	AppendCause() Error
	Causes() []error
	Component() ErrComponent
	SetComponent(c ErrComponent) Error
	SetMessageList(msgList map[int32]string)
}

type MyError struct {
	error
	code        int32
	message     string
	data        map[string]interface{}
	causes      []error
	component   ErrComponent
	appendCause bool
	msgList     map[int32]string
}

func NewMyError(errMsg map[int32]string) Error {
	return &MyError{
		msgList: errMsg,
	}
}

func (e *MyError) Error() string {
	s := fmt.Sprintf("%d:%s", e.code, e.Message())
	if e.appendCause {
		s += getCauses(e.causes)
	}
	return s
}

func (e *MyError) Code() int32 {
	return e.code
}

func (e *MyError) SetCodeMsg(code int32) {
	e.code = code
	e.message = e.msgList[code]
}

func (e *MyError) Message() string {
	return e.message
}

func (e *MyError) SetMessage(msg string) {
	e.message = msg
}

func (e *MyError) SetMessageList(msgList map[int32]string) {
	e.msgList = msgList
}

func (e *MyError) Cause() error {
	if len(e.causes) > 1 {
		return e.causes[0]
	}

	return nil
}

// Causes 用Causes进一步封装，用来保存整个错误堆栈
func (e *MyError) Causes() []error {

	return e.causes
}

// Component 用于识别error发生在哪一层
func (e *MyError) Component() ErrComponent {
	return e.component
}

// SetComponent 设置error发生在哪一层
func (e *MyError) SetComponent(c ErrComponent) Error {
	e.component = c

	return e
}

// getCauses 多个错误封装
func getCauses(errors []error) string {
	var s strings.Builder
	for _, err := range errors {
		s.WriteString(err.Error())
		s.WriteString("; ")
	}
	return s.String()
}

// AppendCause 设置多个cause标识
func (e *MyError) AppendCause() Error {
	e.appendCause = true
	return e
}

func (e *MyError) SetCause(err error) {
	e.causes = append(e.causes, err)
}
