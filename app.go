package main

import (
	"context"
	"fmt"
	"steamflush/internal/steam"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx     context.Context
	scanner *steam.Scanner
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize and resolve Steam directories on startup
	info, err := steam.ResolvePaths()
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to resolve Steam paths: %v", err)
		return
	}

	a.scanner = steam.NewScanner(info)
}

// ScanOrphans acts as a gateway to trigger the isoldated steam service scanner
func (a *App) ScanOrphans() ([]steam.AppInfo, error) {
	if a.scanner == nil {
		return nil, fmt.Errorf("steam scanner service is not initialized")
	}
	return a.scanner.FindOrphans()
}

// DeleteOrphan forwards the deletion request to the internal scanner service
func (a *App) DeleteOrphan(appID string) error {
	if a.scanner == nil {
		return fmt.Errorf("steam scanner service is not initialized")
	}
	return a.scanner.RemoveShaderCache(appID)
}

// RequestConfirmation displays a native operating system message dialog box
func (a *App) RequestConfirmation(title string, message string) (bool, error) {
	result, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         title,
		Message:       message,
		DefaultButton: "No",
	})

	if err != nil {
		return false, err
	}

	return result == "Yes", nil
}
