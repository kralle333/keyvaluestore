package model

import "errors"

var ErrValueNotFound error = errors.New("value not found")
var ErrNoSnapshotsFound error = errors.New("no snapshots found")
