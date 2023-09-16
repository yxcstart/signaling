package common

const (
	NoErr      = 0
	ParamErr   = -1
	NetworkErr = -2
)

type Errors struct {
	code int
	msg  string
}

func New(errCode int, errMsg string) *Errors {
	return &Errors{
		code: errCode,
		msg:  errMsg,
	}
}

func (e *Errors) ErrCode() int {
	return e.code
}

func (e *Errors) ErrorMsg() string {
	return e.msg
}
