package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestFailsGettingTooNew(t *testing.T) {
	actor := NewKeyValueActor(nil, nil, zaptest.NewLogger(t))

	key := "hello"

	actor.putItem(key, "world", 10)
	fetched_value, was_found := actor.getItem(key, 11)

	assert.False(t, was_found)
	assert.Nil(t, fetched_value)
}

func TestCanInsertThenGet(t *testing.T) {
	actor := NewKeyValueActor(nil, nil, zaptest.NewLogger(t))

	var now int64 = 10
	key := "hello"
	value := "world"

	actor.putItem(key, "world", now)
	fetched, was_found := actor.getItem(key, now-1)

	assert.True(t, was_found)
	assert.Equal(t, value, fetched.Value)
}
