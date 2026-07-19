package steam

// AppInfo holds the display information for a game
type AppInfo struct {
	AppID string `json:"appID"`
	Name  string `json:"name"`
	Size  int64  `json:"size"`
}
