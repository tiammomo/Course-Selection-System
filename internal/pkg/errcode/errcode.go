package errcode

// ErrCode 错误码定义
type ErrCode struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

func (e ErrCode) Error() string {
	return e.Msg
}

// 错误码定义
var (
	OK                 = ErrCode{Code: 0, Msg: "success"}
	ParamInvalid       = ErrCode{Code: 1, Msg: "参数不合法"}
	UserHasExisted     = ErrCode{Code: 2, Msg: "该 Username 已存在"}
	UserHasDeleted     = ErrCode{Code: 3, Msg: "用户已删除"}
	UserNotExisted     = ErrCode{Code: 4, Msg: "用户不存在"}
	WrongPassword      = ErrCode{Code: 5, Msg: "密码错误"}
	LoginRequired      = ErrCode{Code: 6, Msg: "用户未登录"}
	CourseNotAvailable = ErrCode{Code: 7, Msg: "课程已满"}
	CourseHasBound     = ErrCode{Code: 8, Msg: "课程已绑定过"}
	CourseNotBind      = ErrCode{Code: 9, Msg: "课程未绑定过"}
	PermDenied         = ErrCode{Code: 10, Msg: "没有操作权限"}
	StudentNotExisted  = ErrCode{Code: 11, Msg: "学生不存在"}
	CourseNotExisted   = ErrCode{Code: 12, Msg: "课程不存在"}
	StudentHasNoCourse = ErrCode{Code: 13, Msg: "学生没有课程"}
	StudentHasCourse   = ErrCode{Code: 14, Msg: "学生有课程"}
	RepeatRequest      = ErrCode{Code: 15, Msg: "重复请求"}
	UnknownError       = ErrCode{Code: 255, Msg: "未知错误"}
)

// WithMsg 创建带消息的错误码
func (e ErrCode) WithMsg(msg string) ErrCode {
	return ErrCode{Code: e.Code, Msg: msg}
}
