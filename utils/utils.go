package utils

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"os"
)

type Response struct {
	Code    int32 `json:"code"`
	Message any   `json:"message"`
	Data    any   `json:"data"`
	Error   any   `json:"error"`
}

func If[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	} else {
		return falseVal
	}
}

func (r *Response) WithData(data any) *Response {
	r.Data = data
	return r
}

func (r *Response) WithError(err any) *Response {
	r.Error = err
	return r
}

func (r *Response) WithMessage(msg any) *Response {
	r.Message = msg
	return r
}

func (r *Response) ToString() string {
	byteStr, _ := json.Marshal(*r)
	return string(byteStr)
}

func (r *Response) Get() Response {
	return *r
}

func (r *Response) ReturnJsonResponse(c *gin.Context, code int) {
	c.JSON(code, r.Get())
}

func NewResponse(code int32, msg string) *Response {
	return &Response{
		Code:    code,
		Message: msg,
		Data:    nil,
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
