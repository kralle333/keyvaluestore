package model

import (
	"errors"
	"time"
)

// REQUESTS AND RESPONSES
type GetValueRequest struct {
	Key         string
	Timestamp   int64
	RespChannel chan GetValueResponse
}

type GetValueResponse struct {
	Value *string
	Error error
}

type PutRequest struct {
	Key       string
	Value     string
	Timestamp int64
}

type SnapshotRequest struct{}

// Used to communicate with the KeyValueActor that is the sole owner of the in memory data
type KeyValueActorCommunication struct {
	Get      chan GetValueRequest
	Put      chan PutRequest
	Snapshot chan SnapshotRequest
}

func NewKeyValueActorCommunication() *KeyValueActorCommunication {
	return &KeyValueActorCommunication{
		Get:      make(chan GetValueRequest),
		Put:      make(chan PutRequest),
		Snapshot: make(chan SnapshotRequest),
	}
}

func (k *KeyValueActorCommunication) GetValue(key string, timestamp int64) (*string, error) {
	receiver := make(chan GetValueResponse)
	k.Get <- GetValueRequest{
		Key:         key,
		Timestamp:   timestamp,
		RespChannel: receiver,
	}
	time := time.NewTimer(time.Duration(time.Duration.Seconds(1)))
	select {
	case <-time.C:
		return nil, errors.New("timed out waiting for Get Value Response")
	case resp := <-receiver:
		return resp.Value, resp.Error
	}
}

func (k *KeyValueActorCommunication) PutValue(key, value string) {
	k.Put <- PutRequest{
		Key:   key,
		Value: value,
	}
}

func (k *KeyValueActorCommunication) TakeSnapshot() {
	k.Snapshot <- SnapshotRequest{}
}
