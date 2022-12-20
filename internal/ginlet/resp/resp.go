package resp

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const DEFAULT_ERROR_CODE = 500
const DEFAULT_ERROR_HTTP_STATUS = 400
const DEFAULT_OK_HTTP_STATUS = 200
const DEFAULT_OK_MSG = "ok"

type StatefulError interface {
	HttpStatusHint() int
	Error() string
}

type ErrorWraper struct {
	error
	statusCode int
}

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type okResp struct {
	Ok any `json:"ok"`
}

type errorResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type errorResp struct {
	Error errorResult `json:"error"`
}

func (e ErrorWraper) HttpStatusHint() int {
	return e.statusCode
}

func NewErrorWraper(err error, code int, wrap string) StatefulError {
	var prefix string
	if code/100 == 4 {
		prefix = "[Client Error]"
	} else if code/100 == 5 {
		prefix = "[Server Interal Error]"
	}
	err = errors.Wrap(err, wrap)
	err = errors.Wrap(err, prefix)
	return ErrorWraper{statusCode: code, error: err}
}

func Error(c *gin.Context, err error) {
	res := errorResp{
		Error: errorResult{
			Code: extractErrorCode(err),
			Msg:  err.Error(),
		},
	}
	c.JSON(extractHttpStatus(err, DEFAULT_ERROR_HTTP_STATUS), res)
}

func ErrorWithHttpStatus(c *gin.Context, err error, httpStatus int) {
	res := errorResp{
		Error: errorResult{
			Code: extractErrorCode(err),
			Msg:  err.Error(),
		},
	}
	c.JSON(httpStatus, res)
}

func Ok(c *gin.Context, result any) {
	if result == nil {
		result = 1
	}
	res := okResp{
		Ok: result,
	}
	c.JSON(extractHttpStatus(result, DEFAULT_OK_HTTP_STATUS), res)
}

func extractErrorCode(a any) int {
	var t, ok = a.(interface {
		ErrorCode() int
	})
	if ok {
		return t.ErrorCode()
	}
	return DEFAULT_ERROR_CODE
}

func extractHttpStatus(a any, d int) int {
	var t, ok = a.(interface {
		HttpStatusHint() int
	})
	if ok {
		return t.HttpStatusHint()
	}
	return d
}
