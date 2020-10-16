package config

import (
	"errors"
	"flag"
	"log"
	"os"
	"strings"
	"syscall"
)

// AppConfig describe configuration for all parts of application
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

// ParseConsole function will parse config options from CLI arguments
func ParseConsole() *AppConfig {
	logPath := flag.String("logPath", "logs/redirective.log", "Path to the log file")

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

	err := checkScreenshotsStorageDir(*screenshotsStoragePath)
	if err != nil {
		log.Fatalf("folder to store screenshots not found and couldn`t be created. error: %s", err)
	}

	if !strings.HasSuffix(*screenshotsStoragePath, "/") {
		*screenshotsStoragePath = *screenshotsStoragePath + "/"
	}


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

func checkScreenshotsStorageDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	_, err := isWritable(path)
	if err != nil {
		return err
	}

	return nil
}

func isWritable(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if !info.IsDir() {
		return false, errors.New("path isn't a directory")
	}

	// Check if the user bit is enabled in file permission
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		return false, errors.New("write permission bit is not set on this file for user")
	}

	var stat syscall.Stat_t
	if err = syscall.Stat(path, &stat); err != nil {
		return false, errors.New("unable to get stat. error " + err.Error())
	}

	if uint32(os.Geteuid()) != stat.Uid {
		return false, errors.New("user doesn't have permission to write to this directory")
	}

	return true, nil
}
