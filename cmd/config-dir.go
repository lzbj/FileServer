package cmd

import (
	"os"
	"path/filepath"
	"sync"

	homedir "github.com/mitchellh/go-homedir"
)

const (
	// Default file server configuration directory where below configuration files/directories are stored.
	// Default file server configuration directory where below configuration files/directories are stored.
	defaultConfigDir = ".fileserver"

	// Default configuration file.
	fileServerConfigFile = "config.json"

	// Directory contains below files/directories for HTTPS configuration.
	certsDir = "certs"

	// Directory contains all CA certificates other than system defaults for HTTPS.
	certsCADir = "CAs"

	// Public certificate file for HTTPS.
	publicCertFile = "public.crt"

	// Private key file for HTTPS.
	privateKeyFile = "private.key"
)

// ConfigDir - configuration directory with locking.
type ConfigDir struct {
	sync.Mutex
	dir string
}

// Set - saves given directory as configuration directory.
func (config *ConfigDir) Set(dir string) {
	config.Lock()
	defer config.Unlock()

	config.dir = dir
}

// Get - returns current configuration directory.
func (config *ConfigDir) Get() string {
	config.Lock()
	defer config.Unlock()

	return config.dir
}

func (config *ConfigDir) getCertsDir() string {
	return filepath.Join(config.Get(), certsDir)
}

// GetCADir - returns certificate CA directory.
func (config *ConfigDir) GetCADir() string {
	return filepath.Join(config.getCertsDir(), certsCADir)
}

// Create - creates configuration directory tree.
func (config *ConfigDir) Create() error {
	return os.MkdirAll(config.GetCADir(), 0700)
}

// GetConfigFile - returns absolute path of config.json file.
func (config *ConfigDir) GetConfigFile() string {
	return filepath.Join(config.Get(), fileServerConfigFile)
}

// GetPublicCertFile - returns absolute path of public.crt file.
func (config *ConfigDir) GetPublicCertFile() string {
	return filepath.Join(config.getCertsDir(), publicCertFile)
}

// GetPrivateKeyFile - returns absolute path of private.key file.
func (config *ConfigDir) GetPrivateKeyFile() string {
	return filepath.Join(config.getCertsDir(), privateKeyFile)
}

func getDefaultConfigDir() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		return ""
	}

	return filepath.Join(homeDir, defaultConfigDir)
}

var configDir = &ConfigDir{dir: getDefaultConfigDir()}

func setConfigDir(dir string) {
	configDir.Set(dir)
}

func getConfigDir() string {
	return configDir.Get()
}

func getCADir() string {
	return configDir.GetCADir()
}

func createConfigDir() error {
	return configDir.Create()
}

func getConfigFile() string {
	return configDir.GetConfigFile()
}

func getPublicCertFile() string {
	return configDir.GetPublicCertFile()
}

func getPrivateKeyFile() string {
	return configDir.GetPrivateKeyFile()
}
