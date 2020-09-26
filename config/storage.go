package config

// StorageConfig describe configuration for storage (database)
type StorageConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Table    string
}
