package model

import "github.com/google/btree"

type KeyValueNodes struct {
	Nodes []KeyValueNode
}

type KeyValueNode struct {
	Key       string
	Value     string
	Timestamp int64
}

func LessThanKeyValueNode(a, b KeyValueNode) bool {
	if a.Key < b.Key {
		return true
	}
	if a.Key > b.Key {
		return false
	}

	return a.Timestamp < b.Timestamp
}

func NewKeyValueTree() *btree.BTreeG[KeyValueNode] {
	return btree.NewG(32, LessThanKeyValueNode)
}
