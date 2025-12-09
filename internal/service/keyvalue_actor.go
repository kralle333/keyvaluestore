package service

import (
	"github.com/google/btree"
	"github.com/kralle333/keyvaluestore/internal/model"
	"github.com/kralle333/keyvaluestore/internal/repository"
	"go.uber.org/zap"
)

// Actor pattern used to have single owner of the key value data, which is stored in a bTree for easy access using third value timestamp
type KeyValueActor struct {
	data            *btree.BTreeG[model.KeyValueNode]
	communication   *model.KeyValueActorCommunication
	storage         *repository.KeyValueStorage
	logger          *zap.Logger
	shutdownChannel chan struct{}
	dirty           bool
}

func NewKeyValueActor(communication *model.KeyValueActorCommunication, storage *repository.KeyValueStorage, parentLogger *zap.Logger) *KeyValueActor {

	return &KeyValueActor{
		data:            model.NewKeyValueTree(),
		communication:   communication,
		storage:         storage,
		logger:          parentLogger.With(zap.String("service", "KeyValueActor")),
		shutdownChannel: make(chan struct{}),
		dirty:           false,
	}
}

func (k *KeyValueActor) PopulateFromSnapshot(snapshot *btree.BTreeG[model.KeyValueNode]) {
	k.logger.Info("Populated data using snapshot")
	k.data = snapshot
}

func (k *KeyValueActor) Spawn() {
	k.logger.Info("Spawning actor communication goroutine")
	go func() {
		for {
			select {
			case <-k.shutdownChannel:
				k.logger.Info("Shutting down")
				return
			case get := <-k.communication.Get:
				k.logger.Info("Getting", zap.String("key", get.Key), zap.Int64("timestamp", get.Timestamp))
				if node, found := k.getItem(get.Key, get.Timestamp); found {
					get.RespChannel <- model.GetValueResponse{
						Value: &node.Value,
						Error: nil,
					}
				} else {
					get.RespChannel <- model.GetValueResponse{
						Value: nil,
						Error: model.ErrValueNotFound,
					}
				}
			case put := <-k.communication.Put:
				k.logger.Debug("Inserting", zap.String("key", put.Key), zap.String("value", put.Value), zap.Int64("timestamp", put.Timestamp))
				k.putItem(put.Key, put.Value, put.Timestamp)
			case <-k.communication.Snapshot:
				if k.dirty {
					copied := k.data.Clone()
					k.storage.SpawnLogSnapshot(copied)
					k.dirty = false
				} else {
					k.logger.Debug("Skipping logging snapshot, dirty=false")
				}
			}
		}
	}()
}

func (k *KeyValueActor) Shutdown() {
	k.shutdownChannel <- struct{}{}
}

func (k *KeyValueActor) putItem(key string, value string, timestamp int64) {
	k.data.ReplaceOrInsert(model.KeyValueNode{
		Key:       key,
		Value:     value,
		Timestamp: timestamp,
	})
	k.dirty = true
}
func (k *KeyValueActor) getItem(key string, timestamp int64) (item *model.KeyValueNode, found bool) {
	pivot := model.KeyValueNode{
		Key:       key,
		Value:     "",
		Timestamp: timestamp,
	}

	var foundItem *model.KeyValueNode
	k.data.AscendGreaterOrEqual(pivot, func(item model.KeyValueNode) bool {
		if item.Key == pivot.Key {
			foundItem = &item
		}
		return foundItem == nil
	})
	return foundItem, foundItem != nil
}
