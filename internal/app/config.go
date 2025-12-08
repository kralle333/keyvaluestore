package app

// Simple struct to hold constants
// TODO: Use viper etc. to load from environment variables and config files
type AppConfig struct {
	SnapshotDir             string
	ListeningPort           uint16
	SnapshotIntervalSeconds int64
}

func GetTestingConfig() AppConfig {
	return AppConfig{
		SnapshotDir:             "test_snapshots/",
		ListeningPort:           8080,
		SnapshotIntervalSeconds: 2,
	}
}
