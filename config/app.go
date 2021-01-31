// Package config provide classes and functions to parse and prepare for use all application options
package config

import (
	"flag"
	"log"
	"os"
	"strings"
	"syscall"
)

// AppConfig describe configuration for all parts of application.
type AppConfig struct {
	AppDomain string
	Storage   *StorageConfig
	// Logger *LogConfig
	HTTPServer      *HTTPServerConfig
	ScreenshotsPath string
	LogsPath        string
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

// ParseConsole function will parse config options from CLI arguments.
func ParseConsole() *AppConfig {
	appDomain := flag.String("appDomain", "redirective.net", "Domain used to host application")
	logPath := flag.String("logPath", "logs", "Path to the log folder")

	screenshotsStoragePath := flag.String("screenshotsPath", "assets/screenshots", "Path to directory where screenshots would be stored")

	host := flag.String("host", "", "Web server listen host")
	port := flag.Int("port", 8080, "Web server listen port")
	certFile := flag.String("certPath", "", "Path to the certificate file")
	keyFile := flag.String("keyPath", "", "Path to the key file")

	storageHost := flag.String("storageHost", "localhost", "Storage host (default: localhost)")
	storagePort := flag.Int("storagePort", 3306, "Storage port (default: 27017)")
	storageUser := flag.String("storageUser", "root", "Storage user (default: root)")
	storagePass := flag.String("storagePass", "secret", "Storage user`s password (default: secret)")
	storageDatabase := flag.String("storageDatabase", "redirective", "Storage database name (default: redirective)")
	storageTable := flag.String("storageTable", "results", "Storage table (default: results)")

	// parse arguments
	flag.Parse()

	err := checkScreenshotsStorageDir(*screenshotsStoragePath)
	if err != nil {
		log.Fatalf("folder to store screenshots not found and couldn`t be created. error: %s", err)
	}

	if !strings.HasSuffix(*screenshotsStoragePath, "/") {
		*screenshotsStoragePath = *screenshotsStoragePath + "/"
	}

	return &AppConfig{
		AppDomain: *appDomain,
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
		LogsPath:        *logPath,
	}
}

// checkScreenshotsStorageDir - check if provided directory exists (or create new) and writeable.
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

// isWritable - check if provided directory is writeable.
func isWritable(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if !info.IsDir() {
		return false, &PathIsNotDirError{}
	}

	// Check if the user bit is enabled in file permission
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		return false, &NoWritePermissionsForUser{}
	}

	var stat syscall.Stat_t
	if err = syscall.Stat(path, &stat); err != nil {
		return false, &UnableToGetStatError{Err: err}
	}

	if uint32(os.Geteuid()) != stat.Uid {
		return false, &NoPermissionsToWriteInDirError{}
	}

	return true, nil
}

// PathIsNotDirError describe error that occurs when wrong path to log dir provided to the AppConfig.
type PathIsNotDirError struct {
}

// Error function return error message.
func (e *PathIsNotDirError) Error() string {
	return `path isn't a directory`
}

// NoWritePermissionsForUser describe error that occurs when application has no permissions to write in log file.
type NoWritePermissionsForUser struct {
}

// Error function return error message.
func (e *NoWritePermissionsForUser) Error() string {
	return `write permission bit is not set on this file for user`
}

// UnableToGetStatError describe error that occurs when application cannot get system information about logs directory.
type UnableToGetStatError struct {
	Err error
}

// Error function return error message.
func (e *UnableToGetStatError) Error() string {
	return "unable to get stat. error: " + e.Err.Error()
}

// NoPermissionsToWriteInDirError describe error that occurs when application has no permissions to write in logs directory.
type NoPermissionsToWriteInDirError struct {
}

// Error function return error message.
func (e *NoPermissionsToWriteInDirError) Error() string {
	return `user doesn't have permission to write to this directory`
}
