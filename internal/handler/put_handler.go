package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kralle333/keyvaluestore/internal/model"
)

type putValueRequest struct {
	Key       string `json:"key" binding:"required"`
	Value     string `json:"value" binding:"required"`
	Timestamp int64  `json:"timestamp" binding:"required"`
}

func (k *KeyValueHttpServer) handlePutValue(c *gin.Context) {
	var req putValueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	k.communication.Put <- model.PutRequest{
		Key:       req.Key,
		Value:     req.Value,
		Timestamp: req.Timestamp,
	}

	c.Status(200)
}
