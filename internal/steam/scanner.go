package steam

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Scanner encapsulates the operations for tracking and cleaning orphan files
type Scanner struct {
	info *Info
}

// NewScanner creates a new scanner service instance using the resolved paths
func NewScanner(info *Info) *Scanner {
	return &Scanner{info: info}
}

// FindOrphans browses the shadercache folder to spot abandoned directories
func (s *Scanner) FindOrphans() ([]AppInfo, error) {
	// Explicit check to feed the front-end empty state workflow
	if _, err := os.Stat(s.info.ShaderCacheDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("shadercache directory not found")
	}

	files, err := os.ReadDir(s.info.ShaderCacheDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read shadercache directory: %w", err)
	}

	var orphans []string
	for _, file := range files {
		if file.IsDir() {
			appID := file.Name()

			// Skip any folders that are not strictly numeric (Steam AppIDs are pure digits)
			if strings.Trim(appID, "0123456789") != "" {
				continue
			}

			// If the appmanifest_[appID].acf file does not exist, the game is uninstalled
			//manifestPath := filepath.Join(s.info.AppsDir, "appmanifest_"+appID+".acf")
			//if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
			orphans = append(orphans, appID)
			//}
		}
	}

	var results []AppInfo
	for _, appID := range orphans {
		name, _ := GetAppName(appID)
		results = append(results, AppInfo{AppID: appID, Name: name})
	}

	return results, nil
}

// RemoveShaderCache safely deletes the target shadercache folder for a specific AppID
func (s *Scanner) RemoveShaderCache(appID string) error {
	if appID == "" {
		return fmt.Errorf("appID cannot be empty")
	}

	targetPath := filepath.Join(s.info.ShaderCacheDir, appID)

	// Security checkpoint: prevent path traversal attacks by ensuring we stay in shadercache
	if !strings.HasPrefix(targetPath, s.info.ShaderCacheDir) {
		return fmt.Errorf("security violation: path traversal detected")
	}

	if _, err := os.Stat(targetPath); err == nil {
		if err := os.RemoveAll(targetPath); err != nil {
			return fmt.Errorf("failed to delete shadercache for %s: %w", appID, err)
		}
	}

	return nil
}
