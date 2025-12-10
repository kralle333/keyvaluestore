package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/google/btree"
	"github.com/kralle333/keyvaluestore/internal/model"
	"go.uber.org/zap"
)

// Responsible for saving snapshots of in memory data to disk
// and can help restore service to last known snapshot
type KeyValueStorage struct {
	dir    string
	logger zap.Logger
}

func NewKeyValueStorage(dir string, parentLogger zap.Logger) *KeyValueStorage {
	return &KeyValueStorage{
		dir:    dir,
		logger: *parentLogger.With(zap.String("source", "keyvalue storage")),
	}
}

// Store the current state and use timestamp to name the file
func (k *KeyValueStorage) SpawnLogSnapshot(tree *btree.BTreeG[model.KeyValueNode], timestamp int64) {

	go func(tree *btree.BTreeG[model.KeyValueNode]) {
		data := model.KeyValueNodes{Nodes: []model.KeyValueNode{}}
		tree.Ascend(func(node model.KeyValueNode) bool {
			data.Nodes = append(data.Nodes, node)
			return true
		})

		jsonData, err := json.Marshal(data)
		if err != nil {
			k.logger.Error("Failed to serialize tree data")
			return
		}

		outputPath := path.Join(k.dir, fmt.Sprintf("state_%d.json", timestamp))
		k.logger.Debug("Attempting to write snapshot to file", zap.String("filepath", outputPath))
		err = os.WriteFile(outputPath, jsonData, os.ModePerm.Perm())
		if err != nil {
			k.logger.Error("Failed to log key value storage to file")
		}
	}(tree)
}

func (k *KeyValueStorage) RetrieveLatest() (*btree.BTreeG[model.KeyValueNode], error) {

	targetFile, err := k.getLatestFile()
	if err != nil {
		return nil, err
	}

	targetPath := path.Join(k.dir, targetFile)
	k.logger.Info("Using found logfile", zap.String("filepath", targetPath))
	data, err := os.ReadFile(targetPath)
	if err != nil {
		return nil, err
	}

	var nodes *model.KeyValueNodes
	err = json.Unmarshal(data, &nodes)
	if err != nil {
		return nil, err
	}

	tree := model.NewKeyValueTree()

	for _, node := range nodes.Nodes {
		k.logger.Info("inserting node", zap.String("key", node.Key), zap.String("value", node.Value), zap.Int64("timestamp", node.Timestamp))
		tree.ReplaceOrInsert(node)
	}

	return tree, nil
}

func (k *KeyValueStorage) getLatestFile() (string, error) {
	entries, err := os.ReadDir(k.dir)
	if err != nil {
		return "", model.ErrUnableToReadFromSnapshotDir
	}

	var jsonFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			jsonFiles = append(jsonFiles, entry.Name())
		}
	}

	if len(jsonFiles) == 0 {
		return "", model.ErrNoSnapshotsFound
	}

	// Extract timestamp and use it to sort the files
	sort.Slice(jsonFiles, func(i, j int) bool {
		getTs := func(name string) int64 {
			base := strings.TrimSuffix(name, ".json")
			part := base[strings.LastIndex(base, "_")+1:]
			ts, _ := strconv.ParseInt(part, 10, 64)
			return ts
		}
		return getTs(jsonFiles[i]) < getTs(jsonFiles[j])
	})

	return jsonFiles[len(jsonFiles)-1], nil
}
