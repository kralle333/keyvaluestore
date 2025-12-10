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
	defer os.RemoveAll(dir)
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
	s.SpawnLogSnapshot(tree, time.Now().Unix())

	// Snapshotting is done in a separate go routine, so wait a bit
	time.Sleep(1 * time.Second)

	latest, err := s.RetrieveLatest()
	assert.NoError(t, err)

	var found *model.KeyValueNode
	latest.AscendGreaterOrEqual(node2, func(item model.KeyValueNode) bool {
		if item.Key == node2.Key {
			found = &item
		}
		return found == nil
	})

	assert.NotNil(t, found)
	assert.Equal(t, node2.Key, found.Key)
	assert.Equal(t, node2.Value, found.Value)
	assert.Equal(t, node2.Timestamp, found.Timestamp)
}

// Sees if parsing the snapshot files results in returning the latest file
func TestCanGetLatest(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	defer os.RemoveAll(dir)
	assert.NoError(t, err)

	s := NewKeyValueStorage(dir, *zaptest.NewLogger(t))

	// Sanity check
	_, err = s.getLatestFile()
	assert.ErrorIs(t, err, model.ErrNoSnapshotsFound)

	tree := model.NewKeyValueTree()
	oldNode := model.KeyValueNode{
		Key:       "thekey",
		Value:     "hey",
		Timestamp: 1,
	}
	tree.ReplaceOrInsert(oldNode)
	s.SpawnLogSnapshot(tree, 10)

	// Snapshotting is done in a separate go routine, so wait a bit
	time.Sleep(1 * time.Second)

	newerTree := model.NewKeyValueTree()
	newNode := model.KeyValueNode{
		Key:       "thekey",
		Value:     "world",
		Timestamp: 1,
	}
	newerTree.ReplaceOrInsert(newNode)
	s.SpawnLogSnapshot(newerTree, 1)

	// Snapshotting is done in a separate go routine, so wait a bit
	time.Sleep(1 * time.Second)

	latest, err := s.RetrieveLatest()
	assert.NoError(t, err)

	var found *model.KeyValueNode
	latest.AscendGreaterOrEqual(model.KeyValueNode{
		Key:       "thekey",
		Value:     "",
		Timestamp: 0,
	}, func(item model.KeyValueNode) bool {
		if item.Key == oldNode.Key {
			found = &item
		}
		return found == nil
	})

	assert.NotNil(t, found)
	assert.Equal(t, oldNode.Key, found.Key)
	assert.Equal(t, oldNode.Value, found.Value)
	assert.Equal(t, oldNode.Timestamp, found.Timestamp)
}
