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

// GetAppNamesBatch retrieves names for a slice of app IDs in parallel.
func GetAppNamesBatch(appIDs []string) (map[string]string, error) {
	cache := loadCache()
	results := make(map[string]string)
	var mu sync.Mutex
	g := new(errgroup.Group)
	g.SetLimit(5)

	for _, id := range appIDs {
		id := id
		if name, ok := cache[id]; ok {
			results[id] = name
			continue
		}

		g.Go(func() error {
			name, err := GetAppName(id)
			if err != nil {
				return err
			}

			mu.Lock()
			results[id] = name
			cache[id] = name
			mu.Unlock()
			return nil
		})
	}

	err := g.Wait()
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
		return "Unknown Game", err
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
		return "Unknown Game", err
	}

	if app, ok := result[appID]; ok && app.Success {
		return app.Data.Name, nil
	}

	return "Unknown Game", nil
}
