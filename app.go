package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"ssh-manager/backend/config"
	"ssh-manager/backend/hosts"
	"ssh-manager/backend/keys"
	"ssh-manager/backend/storage"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx     context.Context
	db      *storage.Database
	hostSvc *hosts.HostService
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	home, _ := os.UserHomeDir()
	dbDir := filepath.Join(home, ".ssh-manager")
	os.MkdirAll(dbDir, 0755)

	db, err := storage.NewDatabase(filepath.Join(dbDir, "data.db"))
	if err != nil {
		runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Database Error",
			Message: fmt.Sprintf("Failed to initialize database: %v", err),
		})
		return
	}
	a.db = db
	a.hostSvc = hosts.NewHostService(db)
}

func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}

// ===== SSH Config Methods =====

func (a *App) GetSSHConfig() map[string]interface{} {
	cfgPath := config.GetDefaultConfigPath()
	cfg, err := config.ParseSSHConfig(cfgPath)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	return map[string]interface{}{
		"filePath":   cfg.FilePath,
		"rawContent": cfg.RawContent,
		"entries":    cfg.Entries,
	}
}

func (a *App) SaveSSHConfig(content string) map[string]interface{} {
	cfgPath := config.GetDefaultConfigPath()
	err := config.SaveSSHConfig(cfgPath, content)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	a.db.AddAuditLog("config_save", "SSH config saved")
	return map[string]interface{}{"success": true}
}

func (a *App) ValidateSSHConfig(content string) map[string]interface{} {
	warnings := config.ValidateSSHConfig(content)
	return map[string]interface{}{"warnings": warnings}
}

func (a *App) ExportConfig(path string) map[string]interface{} {
	content, err := config.ExportConfig(path)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	return map[string]interface{}{"content": content}
}

func (a *App) ImportConfig(path string, content string) map[string]interface{} {
	err := config.ImportConfig(path, content)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	a.db.AddAuditLog("config_import", fmt.Sprintf("Config imported to %s", path))
	return map[string]interface{}{"success": true}
}

// ===== SSH Key Methods =====

func (a *App) ListKeys() map[string]interface{} {
	k, err := keys.ListKeys()
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	return map[string]interface{}{"keys": k}
}

func (a *App) GenerateKey(keyType string, bits int, email string, passphrase string, name string) map[string]interface{} {
	req := keys.KeyGenRequest{
		Type:       keyType,
		Bits:       bits,
		Email:      email,
		Passphrase: passphrase,
		Name:       name,
	}
	path, err := keys.GenerateKey(req)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	a.db.AddAuditLog("key_generate", fmt.Sprintf("Generated %s key: %s", keyType, name))
	return map[string]interface{}{"success": true, "path": path}
}

func (a *App) DeleteKey(name string) map[string]interface{} {
	err := keys.DeleteKey(name)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	a.db.AddAuditLog("key_delete", fmt.Sprintf("Deleted key: %s", name))
	return map[string]interface{}{"success": true}
}

func (a *App) GetPubKeyContent(name string) map[string]interface{} {
	content, err := keys.GetPubKeyContent(name)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	return map[string]interface{}{"content": content}
}

func (a *App) GetPrivKeyContent(name string) map[string]interface{} {
	content, err := keys.GetPrivKeyContent(name)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	return map[string]interface{}{"content": content}
}

// ===== Host Methods =====

func (a *App) ImportHostsFromConfig() map[string]interface{} {
	h, err := a.hostSvc.ImportFromConfig()
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	return map[string]interface{}{"hosts": h, "count": len(h)}
}

func (a *App) GetFavoriteHosts() map[string]interface{} {
	h, err := a.hostSvc.GetFavorites()
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	return map[string]interface{}{"hosts": h}
}

func (a *App) AddFavoriteHost(alias, hostname string, port int, user, tags, group string) map[string]interface{} {
	err := a.hostSvc.AddFavorite(alias, hostname, port, user, tags, group)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	return map[string]interface{}{"success": true}
}

func (a *App) UpdateFavoriteHost(id int64, alias, hostname string, port int, user, tags, group string) map[string]interface{} {
	err := a.hostSvc.UpdateFavorite(id, alias, hostname, port, user, tags, group)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	return map[string]interface{}{"success": true}
}

func (a *App) DeleteFavoriteHost(id int64) map[string]interface{} {
	err := a.hostSvc.DeleteFavorite(id)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	return map[string]interface{}{"success": true}
}

func (a *App) GetGroups() map[string]interface{} {
	g, err := a.hostSvc.GetGroups()
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	return map[string]interface{}{"groups": g}
}

func (a *App) GetConnectionHistory() map[string]interface{} {
	h, err := a.hostSvc.GetHistory()
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	return map[string]interface{}{"history": h}
}

func (a *App) SearchConnectionHistory(query string) map[string]interface{} {
	h, err := a.hostSvc.SearchHistory(query)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	return map[string]interface{}{"history": h}
}

func (a *App) DeleteConnectionHistory(id int64) map[string]interface{} {
	err := a.hostSvc.DeleteHistory(id)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	return map[string]interface{}{"success": true}
}

func (a *App) TestConnection(hostname string, port int) map[string]interface{} {
	result, err := a.hostSvc.TestConnection(hostname, port)
	if err != nil {
		return map[string]interface{}{"error": err.Error(), "success": false}
	}
	return map[string]interface{}{"success": true, "message": result}
}

func (a *App) GetAuditLogs() map[string]interface{} {
	logs, err := a.db.GetAuditLogs()
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	return map[string]interface{}{"logs": logs}
}
