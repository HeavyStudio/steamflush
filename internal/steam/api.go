package steam

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

// ErrNotSteamApp is returned if the ID is not a Steam game.
var ErrNotSteamApp = fmt.Errorf("not a valid steam app")

// GetAppNamesBatch retrieves names for a slice of app IDs in parallel.
func GetAppNamesBatch(appIDs []string) (map[string]string, error) {
	cache := loadCache()
	results := make(map[string]string)

	// fetched holds names resolved by the goroutines. The workers only ever
	// touch this map (guarded by mu) — never `cache` directly — so the
	// unlocked cache reads below stay free of a concurrent map read/write race.
	fetched := make(map[string]string)
	var mu sync.Mutex
	g := new(errgroup.Group)
	g.SetLimit(5)

	for _, id := range appIDs {
		id := id
		if name, ok := cache[id]; ok {
			if name != "" && name != "Unknown Game" {
				results[id] = name
			}
			continue
		}

		g.Go(func() error {
			name, err := GetAppName(id)
			if err != nil {
				if err == ErrNotSteamApp {
					// Cache the negative result so we don't re-query non-apps.
					mu.Lock()
					fetched[id] = ""
					mu.Unlock()
					return nil
				}
				return err
			}

			mu.Lock()
			fetched[id] = name
			mu.Unlock()
			return nil
		})
	}

	err := g.Wait()

	// All goroutines have returned: safe to merge their results and to expose
	// the successfully resolved names to the caller.
	for id, name := range fetched {
		cache[id] = name
		if name != "" {
			results[id] = name
		}
	}

	saveCache(cache)
	return results, err
}

// GetAppName fetches the game name from Steam store API
func GetAppName(appID string) (string, error) {
	safeAppID := url.QueryEscape(appID)
	urlStr := "https://store.steampowered.com/api/appdetails?appids=%s"
	targetURL := fmt.Sprintf(urlStr, safeAppID)

	resp, err := httpClient.Get(targetURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return "", fmt.Errorf("steam API rate limit hit (429)")
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Steam returns a dynamic JSON structure: { "999999": { "success": true, "data": { "name": "..." } } }
	var result map[string]struct {
		Success bool `json:"success"`
		Data    struct {
			Name string `json:"name"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if app, ok := result[appID]; ok && app.Success {
		return app.Data.Name, nil
	}

	return "", ErrNotSteamApp
}
