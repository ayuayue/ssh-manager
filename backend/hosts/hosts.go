package hosts

import (
	"fmt"
	"net"
	"ssh-manager/backend/config"
	"ssh-manager/backend/storage"
	"time"
)

type HostInfo struct {
	Alias    string `json:"alias"`
	HostName string `json:"hostName"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Tags     string `json:"tags"`
	Group    string `json:"group"`
}

type HostService struct {
	db *storage.Database
}

func NewHostService(db *storage.Database) *HostService {
	return &HostService{db: db}
}

func (s *HostService) ImportFromConfig() ([]HostInfo, error) {
	cfgPath := config.GetDefaultConfigPath()
	cfg, err := config.ParseSSHConfig(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	var hosts []HostInfo
	for _, entry := range cfg.Entries {
		port := 22
		if entry.Port != "" {
			fmt.Sscanf(entry.Port, "%d", &port)
		}

		host := HostInfo{
			Alias:    entry.Host,
			HostName: entry.HostName,
			Port:     port,
			User:     entry.User,
			Group:    "default",
		}
		hosts = append(hosts, host)

		existingHosts, _ := s.db.GetFavoriteHosts()
		exists := false
		for _, h := range existingHosts {
			if h.HostAlias == entry.Host {
				exists = true
				break
			}
		}
		if !exists {
			s.db.AddFavoriteHost(entry.Host, entry.HostName, port, entry.User, "", "default")
		}
	}

	return hosts, nil
}

func (s *HostService) GetFavorites() ([]storage.FavoriteHost, error) {
	return s.db.GetFavoriteHosts()
}

func (s *HostService) AddFavorite(alias, hostname string, port int, user, tags, group string) error {
	return s.db.AddFavoriteHost(alias, hostname, port, user, tags, group)
}

func (s *HostService) UpdateFavorite(id int64, alias, hostname string, port int, user, tags, group string) error {
	return s.db.UpdateFavoriteHost(id, alias, hostname, port, user, tags, group)
}

func (s *HostService) DeleteFavorite(id int64) error {
	return s.db.DeleteFavoriteHost(id)
}

func (s *HostService) GetGroups() ([]string, error) {
	return s.db.GetGroups()
}

func (s *HostService) GetHistory() ([]storage.ConnectionHistory, error) {
	return s.db.GetConnectionHistory()
}

func (s *HostService) SearchHistory(query string) ([]storage.ConnectionHistory, error) {
	return s.db.SearchConnectionHistory(query)
}

func (s *HostService) DeleteHistory(id int64) error {
	return s.db.DeleteConnectionHistory(id)
}

func (s *HostService) RecordConnection(alias, hostname string, port int, user string) error {
	return s.db.AddConnectionHistory(alias, hostname, port, user)
}

func (s *HostService) TestConnection(hostname string, port int) (string, error) {
	addr := fmt.Sprintf("%s:%d", hostname, port)

	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return "", fmt.Errorf("connection failed: %s", err.Error())
	}
	defer conn.Close()

	elapsed := time.Since(start)
	return fmt.Sprintf("Connection successful in %v", elapsed), nil
}
