package controller

import (
	"ethbaas/internal/db"
	"ethbaas/internal/server/service"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	dbClient *db.Client
	storeSvc *service.StoreSvc
}

func NewController(dbClient *db.Client) *Controller {
	c := &Controller{
		dbClient: dbClient,
		storeSvc: service.NewStoreSvc(dbClient),
	}
	return c
}

func (c *Controller) Health(ctx *gin.Context) {
	ResponseSuccess(ctx, "Great!")
}
