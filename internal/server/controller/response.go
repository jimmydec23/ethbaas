package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Resp struct {
	Code uint8       `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type PageResp struct {
	Resp
	Total int64 `json:"total"`
}

func ResponseSuccess(ctx *gin.Context, data interface{}) {
	resp := &Resp{
		Code: 0,
		Data: data,
		Msg:  "",
	}
	ctx.JSON(http.StatusOK, resp)
}

func ResponseFail(ctx *gin.Context, msg string) {
	resp := &Resp{
		Code: 1,
		Data: "",
		Msg:  msg,
	}
	ctx.JSON(http.StatusOK, resp)
}

func ResponseLogin(ctx *gin.Context, msg string) {
	resp := &Resp{
		Code: 2,
		Data: "",
		Msg:  msg,
	}
	ctx.JSON(http.StatusOK, resp)
}

func ResponsePage(ctx *gin.Context, data interface{}, total int64) {
	resp := &PageResp{
		Resp: Resp{
			Code: 0,
			Data: data,
			Msg:  "",
		},
		Total: total,
	}
	ctx.JSON(http.StatusOK, resp)
}
