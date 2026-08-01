package main

import (
	"context"
	"errors"
	"fmt"
	"steamflush/internal/steam"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx          context.Context
	scanner      *steam.Scanner
	isSteamFound bool
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize and resolve Steam directories on startup
	info, err := steam.ResolvePaths()
	if err != nil {
		if errors.Is(err, steam.ErrSteamNotFound) {
			a.isSteamFound = false
			runtime.LogWarning(a.ctx, "Steam installation not found")
		} else {
			runtime.LogErrorf(a.ctx, "Fatal error resolving Steam paths: %v", err)
		}
		return
	}

	a.isSteamFound = true
	a.scanner = steam.NewScanner(info)
}

// ScanOrphans acts as a gateway to trigger the isolated steam service scanner
func (a *App) ScanOrphans() ([]steam.AppInfo, error) {
	if !a.isSteamFound {
		return nil, fmt.Errorf("steam installation not detected")
	}

	if a.scanner == nil {
		return nil, fmt.Errorf("steam scanner service is not initialized")
	}
	return a.scanner.FindOrphans()
}

// DeleteOrphan forwards the deletion request to the internal scanner service and records it in history
func (a *App) DeleteOrphan(appID string) error {
	if a.scanner == nil {
		return fmt.Errorf("steam scanner service is not initialized")
	}

	// 1. Récupérer les orphelins pour extraire les métadonnées de celui qu'on va supprimer
	orphans, err := a.scanner.FindOrphans()
	var targetApp *steam.AppInfo
	if err == nil {
		for _, orphan := range orphans {
			if orphan.AppID == appID {
				targetApp = &orphan
				break
			}
		}
	}

	// 2. Supprimer le dossier physique
	if err := a.scanner.RemoveShaderCache(appID); err != nil {
		return err
	}

	// 3. Enregistrer la session dans l'historique
	if targetApp != nil {
		items := []steam.CleanedItem{
			{
				AppID: targetApp.AppID,
				Name:  targetApp.Name,
				Size:  targetApp.Size,
			},
		}

		if saveErr := steam.AddCleanRecord(1, targetApp.Size, items); saveErr != nil {
			runtime.LogErrorf(a.ctx, "Failed to save history record: %v", saveErr)
		}
	}

	return nil
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

func (a *App) IsSteamFound() bool {
	return a.isSteamFound
}

// RemoveShaderCacheBatch handles mass deletion and delegates batch history recording to the scanner
func (a *App) RemoveShaderCacheBatch(appIDs []string) []string {
	if a.scanner == nil {
		return appIDs
	}

	// Le scanner s'occupe déjà de supprimer les dossiers ET de sauvegarder
	// la session dans l'historique via AddCleanRecord
	return a.scanner.RemoveShaderCacheBatch(appIDs)
}

func (a *App) GetCleanHistory() ([]steam.CleanRecord, error) {
	return steam.LoadHistory()
}
