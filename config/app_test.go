package config_test

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/lroman242/redirective/config"
)

func initTestConfig() *config.AppConfig {
	return &config.AppConfig{
		AppDomain: "redirective.net",
		Storage: &config.StorageConfig{
			Host:     "testhost",
			Port:     1234,
			User:     "tester",
			Password: "password",
			Database: "testdb",
			Table:    "testtable",
		},
		HTTPServer: &config.HTTPServerConfig{
			Host:     "test2host",
			Port:     5678,
			HTTPS:    false,
			CertPath: "path/cert.pem",
			KeyPath:  "path/key.pem",
		},
		ScreenshotsPath: "screenshots/",
		LogsPath:        "logs/test.log",
	}
}

func TestParseConsole(t *testing.T) {
	testValues := initTestConfig()
	defer t.Log(removeDir(testValues.ScreenshotsPath))

	appendOsArgs(testValues)

	appConf := config.ParseConsole()

	if appConf.AppDomain != testValues.AppDomain {
		t.Error("invalid value parsed for AppDomain")
	}

	if appConf.LogsPath != testValues.LogsPath {
		t.Error("invalid value parsed for LogsPath")
	}

	if appConf.ScreenshotsPath != testValues.ScreenshotsPath {
		t.Error("invalid value parsed for ScreenshotsPath")
	}

	if appConf.Storage.Host != testValues.Storage.Host {
		t.Error("invalid value parsed for Storage.Host")
	}

	if appConf.Storage.Port != testValues.Storage.Port {
		t.Error("invalid value parsed for Storage.Port")
	}

	if appConf.Storage.User != testValues.Storage.User {
		t.Error("invalid value parsed for Storage.User")
	}

	if appConf.Storage.Password != testValues.Storage.Password {
		t.Error("invalid value parsed for Storage.Password")
	}

	if appConf.Storage.Database != testValues.Storage.Database {
		t.Error("invalid value parsed for Storage.Database")
	}

	if appConf.Storage.Table != testValues.Storage.Table {
		t.Error("invalid value parsed for Storage.Table")
	}

	if appConf.HTTPServer.Host != testValues.HTTPServer.Host {
		t.Error("invalid value parsed for HTTPServer.Host")
	}

	if appConf.HTTPServer.Port != testValues.HTTPServer.Port {
		t.Error("invalid value parsed for HTTPServer.Port")
	}

	if appConf.HTTPServer.CertPath != testValues.HTTPServer.CertPath {
		t.Error("invalid value parsed for HTTPServer.CertPath")
	}

	if appConf.HTTPServer.KeyPath != testValues.HTTPServer.KeyPath {
		t.Error("invalid value parsed for HTTPServer.KeyPath")
	}
}

func removeDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}

	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}

	err = os.Remove(dir)
	if err != nil {
		return err
	}

	return nil
}

func appendOsArgs(testValues *config.AppConfig) {
	os.Args = append(os.Args, "--appDomain="+testValues.AppDomain)
	os.Args = append(os.Args, "--logPath="+testValues.LogsPath)
	os.Args = append(os.Args, "--screenshotsPath="+testValues.ScreenshotsPath)
	os.Args = append(os.Args, "--host="+testValues.HTTPServer.Host)
	os.Args = append(os.Args, "--port="+strconv.Itoa(testValues.HTTPServer.Port))
	os.Args = append(os.Args, "--certPath="+testValues.HTTPServer.CertPath)
	os.Args = append(os.Args, "--keyPath="+testValues.HTTPServer.KeyPath)
	os.Args = append(os.Args, "--storageHost="+testValues.Storage.Host)
	os.Args = append(os.Args, "--storagePort="+strconv.Itoa(testValues.Storage.Port))
	os.Args = append(os.Args, "--storageUser="+testValues.Storage.User)
	os.Args = append(os.Args, "--storagePass="+testValues.Storage.Password)
	os.Args = append(os.Args, "--storageDatabase="+testValues.Storage.Database)
	os.Args = append(os.Args, "--storageTable="+testValues.Storage.Table)
}
