package steam

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type CleanedItem struct {
	AppID string `json:"app_id"`
	Name  string `json:"name"`
	Size  int64  `json:"size"`
}

type CleanRecord struct {
	Timestamp  string        `json:"timestamp"`
	ItemsCount int           `json:"items_count"`
	BytesFreed int64         `json:"bytes_freed"`
	Items      []CleanedItem `json:"items"`
}

type HistoryStore struct {
	Records []CleanRecord `json:"records"`
}

// Get the history file
func getHistoryFilePath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	appDir := filepath.Join(configDir, "steamflush")
	os.MkdirAll(appDir, 0755)
	return filepath.Join(appDir, "history.json")
}

// Load history from disk
func LoadHistory() ([]CleanRecord, error) {
	filepath := getHistoryFilePath()
	data, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return []CleanRecord{}, nil
		}
		return nil, err
	}

	var store HistoryStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}

	return store.Records, nil
}

// Save a new cleaning session
func AddCleanRecord(itemsCount int, bytesFreed int64, items []CleanedItem) error {
	records, err := LoadHistory()
	if err != nil {
		records = []CleanRecord{}
	}

	newRecord := CleanRecord{
		Timestamp:  time.Now().Format(time.RFC3339),
		ItemsCount: itemsCount,
		BytesFreed: bytesFreed,
		Items:      items,
	}

	records = append([]CleanRecord{newRecord}, records...)

	store := HistoryStore{Records: records}
	data, err := json.MarshalIndent(store, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(getHistoryFilePath(), data, 0644)
}
