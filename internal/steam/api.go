package steam

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetAppName fetches the game name from Steam store API
func GetAppName(appID string) (string, error) {
	url := fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=%s", appID)
	resp, err := http.Get(url)
	if err != nil {
		return "Unknown App", err
	}
	defer resp.Body.Close()

	// Steam returns a dynamic JSON structure: { "999999": { "success": true, "data": { "name": "..." } } }
	var result map[string]struct {
		Success bool `json:"success"`
		Data    struct {
			Name string `json:"name"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "Unknown App", err
	}

	if app, ok := result[appID]; ok && app.Success {
		return app.Data.Name, nil
	}

	return "Unknown App", nil
}
