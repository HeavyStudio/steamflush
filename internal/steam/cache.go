package steam

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func getCachePath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cacheDir, "steamflush", "names.json"), nil
}

func loadCache() map[string]string {
	path, err := getCachePath()
	if err != nil {
		return make(map[string]string)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return make(map[string]string)
	}

	var cache map[string]string
	if err := json.Unmarshal(data, &cache); err != nil {
		fmt.Printf("Warning: failed to unmarshal cache, resetting: %v\n", err)
		return make(map[string]string)
	}
	return cache
}

func saveCache(cache map[string]string) {
	path, err := getCachePath()
	if err != nil {
		fmt.Printf("Warning: failed to get cache path: %v\n", err)
		return
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		fmt.Printf("Warning: failed to create cache directory: %v\n", err)
		return
	}

	data, err := json.Marshal(cache)
	if err != nil {
		fmt.Printf("Warning: failed to marshal cache: %v\n", err)
		return
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		fmt.Printf("Warning: failed to save cache to disk: %v\n", err)
	}
}
