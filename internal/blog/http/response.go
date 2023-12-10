package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	Success             = 0
	Failed              = 1
	Exception           = 1001
	ParametersIncorrect = 1002
	NotSupportedEvent   = 1004

	GatewayNotFound = 404
	BackendNotLogin = 419
)

var ResponseMessages = map[int]string{
	Success:             "成功",
	Failed:              "失败",
	Exception:           "异常",
	ParametersIncorrect: "参数错误",
	NotSupportedEvent:   "事件不存在",
	GatewayNotFound:     "借口不存在",
}

func ShowSuccessWithData(c *gin.Context, data interface{}) {
	var rsp interface{}
	rsp = gin.H{"err_code": Success, "message": "success", "data": data}

	c.JSON(http.StatusOK, rsp)
	c.Abort()
	return
}

func ShowSuccess(c *gin.Context) {
	var rsp interface{}
	rsp = gin.H{
		"err_code": Success,
		"message":  "success",
	}

	c.JSON(http.StatusOK, rsp)
	c.Abort()
	return
}

func ShowError(c *gin.Context, message string) {
	rsp := gin.H{
		"err_code": Failed,
		"message":  message,
	}
	c.JSON(http.StatusOK, rsp)
	c.Abort()
	return
}

func ShowErrorCode(c *gin.Context, code int) {
	rsp := gin.H{"err_code": code, "message": ResponseMessages[code]}
	c.JSON(http.StatusOK, rsp)
	c.Abort()
	return
}

func ShowException(c *gin.Context, traceNo string) {
	msg := "something wrong, please contact out customer service with trace no. #%s"
	msg = fmt.Sprintf(msg, traceNo)
	rsp := gin.H{"err_code": Exception, "message": msg}
	c.JSON(http.StatusOK, rsp)
	c.Abort()
	return
}
