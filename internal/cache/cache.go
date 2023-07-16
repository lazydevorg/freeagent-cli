package cache

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
)

func SaveJson(name string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	base64Data := base64.StdEncoding.EncodeToString(jsonData)
	return Save(name, []byte(base64Data))
}

func LoadJson(name string, data interface{}) error {
	base64Data, err := Load(name)
	if err != nil {
		return err
	}
	jsonData, err := base64.StdEncoding.DecodeString(string(base64Data))
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, data)
}

func Save(name string, data []byte) error {
	path, err := getFilePath(name)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func Load(name string) ([]byte, error) {
	path, err := getFilePath(name)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}

func getCachePath() (string, error) {
	cachePath := os.Getenv("CACHE_PATH")
	if cachePath != "" {
		return cachePath, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	cachePath = filepath.Join(home, ".cache", "freeagent-cli")
	err = os.MkdirAll(cachePath, 0700)
	if err != nil {
		return "", err
	}
	return cachePath, nil
}

func getFilePath(name string) (string, error) {
	cachePath, err := getCachePath()
	if err != nil {
		return "", err
	}
	return filepath.Join(cachePath, name), nil
}
