package controller

import (
	"ethbaas/internal/server/message"

	"github.com/gin-gonic/gin"
)

func (c *Controller) StoreQuery(ctx *gin.Context) {
	req := &message.StoreQuery{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	value, err := c.storeSvc.Query(req.Key)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	ResponseSuccess(ctx, value)
}

func (c *Controller) StoreWrite(ctx *gin.Context) {
	req := &message.StoreWrite{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	tx, err := c.storeSvc.Write(req.Key, req.Value)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	ResponseSuccess(ctx, tx)
}
