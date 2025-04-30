package gerr

import "fmt"

var errMap = make(map[int32]*Error)

func Register(err ...*Error) {
	for _, gerr := range err {
		errMap[gerr.code] = gerr
	}
}

func GetErrCode(err error) int32 {
	if err == nil {
		return 0
	}
	if p, ok := err.(*Error); ok {
		return p.code
	}

	return -1
}

func GetErrMsg(errCode int32) string {
	if errCode == 0 {
		return "success"
	}
	msg, ok := errMap[errCode]
	if ok {
		return msg.message
	}
	if errCode < 0 {
		return "system error"
	}
	return "unknown"
}

func GetErrMsgByErr(err error) string {
	if x, ok := err.(*Error); ok {
		return x.message
	} else {
		return err.Error()
	}
}

func Creategerr(code int32) *Error {
	gerr, ok := errMap[code]
	if ok {
		return gerr
	}
	return &Error{code: code, message: fmt.Sprintf("unknown code %d", code)}
}
