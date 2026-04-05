package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SSHHostEntry struct {
	Host          string `json:"host"`
	HostName      string `json:"hostName"`
	User          string `json:"user"`
	Port          string `json:"port"`
	IdentityFile  string `json:"identityFile"`
	ForwardAgent  string `json:"forwardAgent"`
	ProxyJump     string `json:"proxyJump"`
	ExtraOptions  []string `json:"extraOptions"`
	RawBlock      string `json:"rawBlock"`
}

type SSHConfig struct {
	FilePath string         `json:"filePath"`
	Entries  []SSHHostEntry `json:"entries"`
	RawContent string       `json:"rawContent"`
}

func GetDefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".ssh", "config")
}

func ParseSSHConfig(path string) (*SSHConfig, error) {
	if path == "" {
		path = GetDefaultConfigPath()
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := &SSHConfig{
		FilePath:   path,
		RawContent: string(content),
	}

	entries := parseSSHConfigContent(string(content))
	cfg.Entries = entries

	return cfg, nil
}

func parseSSHConfigContent(content string) []SSHHostEntry {
	var entries []SSHHostEntry
	var currentHost *SSHHostEntry

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		parts := strings.SplitN(trimmed, " ", 2)
		if len(parts) < 2 {
			continue
		}

		key := strings.ToLower(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, "\"")

		if key == "host" {
			if currentHost != nil {
				entries = append(entries, *currentHost)
			}
			currentHost = &SSHHostEntry{
				Host:         value,
				Port:         "22",
				ExtraOptions: []string{},
			}
		} else if currentHost != nil {
			switch key {
			case "hostname":
				currentHost.HostName = value
			case "user":
				currentHost.User = value
			case "port":
				currentHost.Port = value
			case "identityfile":
				currentHost.IdentityFile = value
			case "forwardagent":
				currentHost.ForwardAgent = value
			case "proxyjump":
				currentHost.ProxyJump = value
			default:
				currentHost.ExtraOptions = append(currentHost.ExtraOptions, trimmed)
			}
		}
	}

	if currentHost != nil {
		entries = append(entries, *currentHost)
	}

	return entries
}

func SaveSSHConfig(path string, content string) error {
	if path == "" {
		path = GetDefaultConfigPath()
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func ValidateSSHConfig(content string) []string {
	var warnings []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			warnings = append(warnings, fmt.Sprintf("Line %d: Invalid format '%s'", lineNum, line))
			continue
		}

		key := strings.ToLower(parts[0])
		validKeys := map[string]bool{
			"host": true, "hostname": true, "user": true, "port": true,
			"identityfile": true, "forwardagent": true, "proxyjump": true,
			"proxycommand": true, "stricthostkeychecking": true,
			"userknownhostsfile": true, "identitiesonly": true,
			"addkeystoagent": true, "serveraliveinterval": true,
			"serveralivecountmax": true, "loglevel": true,
			"compression": true, "ciphers": true, "macs": true,
			"kexalgorithms": true, "hostkeyalgorithms": true,
			"pubkeyauthentication": true, "passwordauthentication": true,
			"batchmode": true, "connecttimeout": true,
		}

		if !validKeys[key] {
			warnings = append(warnings, fmt.Sprintf("Line %d: Unknown keyword '%s'", lineNum, key))
		}

		if key == "port" {
			if parts[1] != "22" && parts[1] != "2222" {
				port := strings.TrimSpace(parts[1])
				if port != "22" {
					// Just note non-standard ports
				}
			}
		}
	}

	if len(warnings) == 0 {
		warnings = append(warnings, "Config is valid")
	}

	return warnings
}

func ExportConfig(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func ImportConfig(path string, content string) error {
	return SaveSSHConfig(path, content)
}
