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
	info, err := os.Stat(s.info.ShaderCacheDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("shadercache directory not found: %w", err)
		}
		return nil, fmt.Errorf("failed to access shadercache directory: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("shadercache path is not a directory")
	}

	files, err := os.ReadDir(s.info.ShaderCacheDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read shadercache directory: %w", err)
	}

	var orphans []string
	for _, file := range files {
		if file.IsDir() {
			appID := file.Name()

			if strings.Trim(appID, "0123456789") != "" {
				continue
			}

			manifestPath := filepath.Join(s.info.AppsDir, "appmanifest_"+appID+".acf")
			_, err := os.Stat(manifestPath)
			if err != nil {
				if os.IsNotExist(err) {
					orphans = append(orphans, appID)
				} else {
					fmt.Printf("Warning: failed to check manifest for app %s: %v\n", appID, err)
				}
			}
		}
	}

	names, err := GetAppNamesBatch(orphans)
	if err != nil {
		fmt.Printf("Warning: failed to fetch game names: %v\n", err)
	}

	var results []AppInfo
	for _, appID := range orphans {
		name, ok := names[appID]
		if !ok || name == "" {
			continue
		}

		folderPath := filepath.Join(s.info.ShaderCacheDir, appID)
		size := s.GetDirectorySize(folderPath)

		results = append(results, AppInfo{
			AppID: appID,
			Name:  name,
			Size:  size,
		})
	}

	return results, nil
}

// GetOrphanInfo resolves the display metadata (name + size) for a single
// shadercache entry without scanning the whole directory. The size is read
// before any deletion; the name is resolved cache-first, so in the normal
// flow (the name was already cached during the scan) this triggers no network
// request. Falls back to a placeholder name if resolution fails.
func (s *Scanner) GetOrphanInfo(appID string) (AppInfo, error) {
	if appID == "" || strings.Trim(appID, "0123456789") != "" {
		return AppInfo{}, fmt.Errorf("invalid appID")
	}

	folderPath := filepath.Join(s.info.ShaderCacheDir, appID)
	size := s.GetDirectorySize(folderPath)

	names, err := GetAppNamesBatch([]string{appID})
	if err != nil {
		// Non-fatal: keep the size and fall back to a placeholder name.
		fmt.Printf("Warning: failed to fetch name for app %s: %v\n", appID, err)
	}

	name := names[appID]
	if name == "" {
		name = "Unknown Game (" + appID + ")"
	}

	return AppInfo{AppID: appID, Name: name, Size: size}, nil
}

// GetDirectorySize calculates the total size of a directory by walking its files
func (s *Scanner) GetDirectorySize(path string) int64 {
	var size int64
	err := filepath.WalkDir(path, func(currentPath string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err == nil {
				size += info.Size()
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Warning: failed to fully calculate size for %s: %v\n", path, err)
	}

	return size
}

// RemoveShaderCache safely deletes the target shadercache folder for a specific AppID
func (s *Scanner) RemoveShaderCache(appID string) error {
	if appID == "" || strings.Trim(appID, "0123456789") != "" {
		return fmt.Errorf("invalid appID")
	}

	targetPath := filepath.Join(s.info.ShaderCacheDir, appID)
	cleaned := filepath.Clean(targetPath)

	if !strings.HasPrefix(cleaned, s.info.ShaderCacheDir+string(filepath.Separator)) {
		return fmt.Errorf("security violation: path traversal detected")
	}

	if err := os.RemoveAll(cleaned); err != nil {
		return fmt.Errorf("failed to delete shadercache for %s: %w", appID, err)
	}
	return nil
}

// RemoveShaderCacheBatch iterates through a list of AppIDs and attempts to remove their shadercache folders
func (s *Scanner) RemoveShaderCacheBatch(appIDs []string) []string {
	names, err := GetAppNamesBatch(appIDs)
	if err != nil {
		fmt.Printf("Warning: failed to fetch game names for history: %v\n", err)
		names = make(map[string]string)
	}

	var failedIDs []string
	var cleanedItems []string
	var historyItems []CleanedItem
	var totalBytesFreed int64

	for _, id := range appIDs {
		dirPath := filepath.Join(s.info.ShaderCacheDir, id)
		size := s.GetDirectorySize(dirPath)

		name := names[id]
		if name == "" {
			name = "Unknown Game (" + id + ")"
		}

		if err := s.RemoveShaderCache(id); err != nil {
			failedIDs = append(failedIDs, id)
		} else {
			cleanedItems = append(cleanedItems, id)
			totalBytesFreed += size
			historyItems = append(historyItems, CleanedItem{
				AppID: id,
				Name:  name,
				Size:  size,
			})
		}
	}

	if len(historyItems) > 0 {
		if err := AddCleanRecord(len(historyItems), totalBytesFreed, historyItems); err != nil {
			fmt.Printf("Warning: failed to save cleaning history: %v\n", err)
		}
	}

	return failedIDs
}
