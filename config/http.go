package config

// HTTPServerConfig describe configuration for web server.
type HTTPServerConfig struct {
	Host     string
	Port     int
	HTTPS    bool
	CertPath string
	KeyPath  string
}
