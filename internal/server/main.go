package server

import (
	"ethbaas/internal/config"
	"ethbaas/internal/db"
	"ethbaas/internal/log"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"ethbaas/internal/server/controller"

	"github.com/gin-gonic/gin"
)

type Server struct{}

func NewServer() *Server {
	s := &Server{}
	return s
}

func (s *Server) Start() {
	dbClient, err := db.NewClient()
	if err != nil {
		panic(err)
	}
	defer dbClient.Close()

	c := controller.NewController(dbClient)
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", c.Health)
		v1.POST("/store/query", c.StoreQuery)
		v1.POST("/store/write", c.StoreWrite)
	}

	port := config.C.GetInt("server.port")
	router.Run(fmt.Sprintf(":%d", port))
	log.Logger.Info("Logic server started.")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	log.Logger.Info("Server shutdown.")
}
