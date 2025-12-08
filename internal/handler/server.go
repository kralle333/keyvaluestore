package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kralle333/keyvaluestore/internal/model"
)

type KeyValueHttpServer struct {
	router        *gin.Engine
	communication *model.KeyValueActorCommunication
	listeningPort uint16
}

func NewKeyValueHttpServer(communication *model.KeyValueActorCommunication, listeningPort uint16) *KeyValueHttpServer {

	r := gin.New()
	r.Use(gin.Recovery())
	// Just use default logger for nice coloring etc.
	r.Use(gin.Logger())

	return &KeyValueHttpServer{
		router:        r,
		communication: communication,
		listeningPort: listeningPort,
	}
}

func (k *KeyValueHttpServer) ListenAndServe() error {

	// Register routes
	k.router.GET("/", k.handleGetValue)
	k.router.PUT("/", k.handlePutValue)
	k.router.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	addr := fmt.Sprintf(":%d", k.listeningPort)
	return k.router.Run(addr)
}
