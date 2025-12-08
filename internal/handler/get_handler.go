package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kralle333/keyvaluestore/internal/model"
)

type getValueRequest struct {
	Key       string `json:"key" binding:"required"`
	Timestamp int64  `json:"timestamp" binding:"required"`
}

type getValueResponse struct {
	Value string `json:"value"`
}

func (k *KeyValueHttpServer) handleGetValue(c *gin.Context) {
	var req getValueRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	respChan := make(chan model.GetValueResponse)
	k.communication.Get <- model.GetValueRequest{
		Key:         req.Key,
		Timestamp:   req.Timestamp,
		RespChannel: respChan,
	}

	resp := <-respChan
	if resp.Error != nil {
		if errors.Is(resp.Error, model.ErrValueNotFound) {
			c.AbortWithStatus(404)
		} else {
			c.AbortWithStatus(500)
		}
		return
	}

	c.IndentedJSON(200, getValueResponse{
		Value: *resp.Value,
	})

}
