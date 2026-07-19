package steam

import (
	"fmt"
	"os"
	"path/filepath"
)

// Info holds the strategic directory paths for Steam storage management
type Info struct {
	AppsDir        string
	ShaderCacheDir string
	CompatDataDir  string
}

// ResolvePaths automatically detects the actual Steam installation directories on Linux
func ResolvePaths() (*Info, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	// List of standard Steam installation paths on Linux/Debian to probe sequentially
	possiblePaths := []string{
		filepath.Join(home, ".steam", "debian-installation"),                                      // Standard Debian apt package
		filepath.Join(home, ".local", "share", "Steam"),                                           // Standard Valve/Ubuntu path
		filepath.Join(home, ".steam", "steam"),                                                    // Legacy/Arch symlink path
		filepath.Join(home, ".var", "app", "com.valvesoftware.Steam", ".local", "share", "Steam"), // Flatpak sandbox
	}

	var steamBase string
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			steamBase = path
			break
		}
	}

	// If no directory matches, fallback to the standard location layout
	if steamBase == "" {
		steamBase = filepath.Join(home, ".steam", "debian-installation")
	}

	appsDir := filepath.Join(steamBase, "steamapps")

	return &Info{
		AppsDir:        appsDir,
		ShaderCacheDir: filepath.Join(appsDir, "shadercache"),
		CompatDataDir:  filepath.Join(appsDir, "compatdata"),
	}, nil
}
