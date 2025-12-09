package repository

import (
	"os"
	"testing"
	"time"

	"github.com/kralle333/keyvaluestore/internal/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

// Check that we can store btree data and then read it from file again
func TestCanStoreAndRetrieve(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	// defer os.RemoveAll(dir)
	assert.NoError(t, err)

	s := NewKeyValueStorage(dir, *zaptest.NewLogger(t))

	// Sanity check
	_, err = s.getLatestFile()
	assert.ErrorIs(t, err, model.ErrNoSnapshotsFound)

	tree := model.NewKeyValueTree()

	node1 := model.KeyValueNode{
		Key:       "hello",
		Value:     "world",
		Timestamp: 1,
	}
	node2 := model.KeyValueNode{
		Key:       "cool",
		Value:     "world",
		Timestamp: 10,
	}

	tree.ReplaceOrInsert(node1)
	tree.ReplaceOrInsert(node2)
	s.SpawnLogSnapshot(tree)

	// Snapshotting is done in a separate go routine, so wait a bit
	time.Sleep(1 * time.Second)

	latest, err := s.RetrieveLatest()
	assert.NoError(t, err)

	var found *model.KeyValueNode
	latest.AscendGreaterOrEqual(node1, func(item model.KeyValueNode) bool {
		if item.Key == node1.Key {
			found = &item
		}
		return found == nil
	})

	assert.NotNil(t, found)
	assert.Equal(t, node1.Key, found.Key)
	assert.Equal(t, node1.Value, found.Value)
	assert.Equal(t, node1.Timestamp, found.Timestamp)
}
