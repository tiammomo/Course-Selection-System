package response

import (
	"net/http"

	"course_select/internal/pkg/errcode"
)

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"message,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(data interface{}) Response {
	return Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
}

// Fail 失败响应
func Fail(code errcode.ErrCode) Response {
	return Response{
		Code: code.Code,
		Msg:  code.Msg,
	}
}

// FailWithMsg 带消息的失败响应
func FailWithMsg(code errcode.ErrCode, msg string) Response {
	return Response{
		Code: code.Code,
		Msg:  msg,
	}
}

// FailWithError 带错误的失败响应
func FailWithError(err error) Response {
	if e, ok := err.(errcode.ErrCode); ok {
		return Fail(e)
	}
	return Fail(errcode.UnknownError)
}

// Unauthorized 未授权
func Unauthorized(msg string) Response {
	return Response{
		Code: http.StatusUnauthorized,
		Msg:  msg,
	}
}

// Forbidden 禁止访问
func Forbidden(msg string) Response {
	return Response{
		Code: http.StatusForbidden,
		Msg:  msg,
	}
}
