# Simple Key Value Store
- Supports Putting and Getting via HTTP
- Automatic snaphotting to disk
- Automatic restore using latest snapshot on startup

## Configuration
Currently all configuration values are hardcoded with the following options:
- Snapshot interval: 2s
- Restore from backup: true
- Listening port: 8080
- Snapshot dir: test_snapshots/
## Instructions
- Run using `go run cmd/api/main.go`
- A simple test can be run using `make curl_test` that puts and the gets afterwards
- Run all tests using `go test ./... -v`
