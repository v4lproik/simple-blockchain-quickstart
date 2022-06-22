package domains

import (
	"encoding/json"
	"net/http"
)

type ErrorBuilder struct{}

func NewErrorBuilder() ErrorBuilder {
	return ErrorBuilder{}
}

func (e *ErrorBuilder) New(code int, msg string, context ...interface{}) *Error {
	resp := e.NewUnknownError()

	resp.Err.Code = code
	resp.Err.Status = http.StatusText(code)
	resp.Err.Message = msg
	resp.Err.Context = make([]interface{}, 0)

	if code < 1 {
		resp.Err.Code = http.StatusInternalServerError
		resp.Err.Status = http.StatusText(http.StatusInternalServerError)
	}

	for i := range context {
		c := context[i]
		if v, ok := c.(error); ok {
			resp.Err.Context = append(resp.Err.Context, v.Error())
		} else {
			resp.Err.Context = append(resp.Err.Context, c)
		}
	}

	return resp
}

func (e *ErrorBuilder) NewUnknownError() *Error {
	return e.New(http.StatusInternalServerError, "")
}

type Error struct {
	Err struct {
		Code    int           `json:"code"`
		Status  string        `json:"status"`
		Message string        `json:"message"`
		Context []interface{} `json:"context"`
	} `json:"error"`
}

func (err *Error) Code() int {
	return err.Err.Code
}

func (err *Error) Message() string {
	return err.Err.Message
}

func (err *Error) Context() []interface{} {
	return err.Err.Context
}

func (err *Error) Error() string {
	res := "{}"
	if bytes, marshalErr := json.Marshal(err); marshalErr == nil {
		res = string(bytes)
	}

	return res
}
