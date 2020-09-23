package config

import "flag"

type AppConfig struct {
	Storage *StorageConfig
	//Logger *LogConfig
	HTTPServer      *HTTPServerConfig
	ScreenshotsPath string
	LogFilePath     string
}

//
//func ParseYAML(path string) *AppConfig {
//	//TODO: parse from yaml
//}
//
//func ParseJSON(path string) *AppConfig {
//	// TODO: parse from json
//}
//
//func ParseENV(path string) *AppConfig {
////	// TODO: parse from .env
//}
//

func ParseConsole() *AppConfig {
	logPath := flag.String("logPath", "log/redirective.log", "Path to the log file")

	screenshotsStoragePath := flag.String("screenshotsPath", "assets/screenshots", "Path to directory where screenshots would be stored")

	host := flag.String("host", "", "Web server listen host")
	port := flag.Int("port", 8080, "Web server listen port")
	certFile := flag.String("certPath", "/etc/ssl/cert.pem", "Path to the certificate file")
	keyFile := flag.String("keyPath", "/etc/ssl/privkey.pem", "Path to the key file")

	storageHost := flag.String("storageHost", "localhost", "Storage host (default: localhost)")
	storagePort := flag.Int("storagePort", 3306, "Storage port (default: 27017)")
	storageUser := flag.String("storageUser", "root", "Storage user (default: root)")
	storagePass := flag.String("storagePass", "secret", "Storage user`s password (default: secret)")
	storageDatabase := flag.String("storageDatabase", "redirective", "Storage database name (default: redirective)")
	storageTable := flag.String("storageTable", "results", "Storage table (default: results)")

	//parse arguments
	flag.Parse()

	return &AppConfig{
		Storage: &StorageConfig{
			Host:     *storageHost,
			Port:     *storagePort,
			User:     *storageUser,
			Password: *storagePass,
			Database: *storageDatabase,
			Table:    *storageTable,
		},
		HTTPServer: &HTTPServerConfig{
			Host:     *host,
			Port:     *port,
			HTTPS:    *certFile != "" && *keyFile != "",
			CertPath: *certFile,
			KeyPath:  *keyFile,
		},
		ScreenshotsPath: *screenshotsStoragePath,
		LogFilePath:     *logPath,
	}
}
