package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func GetLocalDataDir() (string, error) {
	appName := "flashback"

	var dataDir string
	var err error

	switch runtime.GOOS {
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dataDir = filepath.Join(homeDir, "Library", "Application Support")
	case "linux":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dataDir = filepath.Join(homeDir, ".local", "share")
	case "windows":
		dataDir = os.Getenv("LOCALAPPDATA")
		if dataDir == "" {
			return "", fmt.Errorf("error: %%LocalAppData%% environment variable not found")
		}
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	appDataDir := filepath.Join(dataDir, appName)

	err = os.MkdirAll(appDataDir, 0755)
	if err != nil {
		return "", fmt.Errorf("error creating app data directory: %w", err)
	}

	return appDataDir, nil
}

func GetConfigDir() (string, error) {
	appName := "flashback"

	var configDir string
	var err error

	configDir, err = os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appConfigDir := filepath.Join(configDir, appName)

	err = os.MkdirAll(appConfigDir, 0755)
	if err != nil {
		return "", fmt.Errorf("error creating app data directory: %w", err)
	}

	return appConfigDir, nil
}
