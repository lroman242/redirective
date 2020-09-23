package config

type HTTPServerConfig struct {
	Host     string
	Port     int
	HTTPS    bool
	CertPath string
	KeyPath  string
}
