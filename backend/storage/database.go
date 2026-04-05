package storage

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ConnectionHistory struct {
	ID           int64     `json:"id"`
	HostAlias    string    `json:"hostAlias"`
	HostName     string    `json:"hostName"`
	Port         int       `json:"port"`
	User         string    `json:"user"`
	LastConnected time.Time `json:"lastConnected"`
}

type FavoriteHost struct {
	ID        int64  `json:"id"`
	HostAlias string `json:"hostAlias"`
	HostName  string `json:"hostName"`
	Port      int    `json:"port"`
	User      string `json:"user"`
	Tags      string `json:"tags"`
	Group     string `json:"group"`
}

type AuditLog struct {
	ID        int64     `json:"id"`
	Action    string    `json:"action"`
	Detail    string    `json:"detail"`
	Timestamp time.Time `json:"timestamp"`
}

type Database struct {
	db *sql.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	d := &Database{db: db}
	if err := d.initTables(); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Database) initTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS connection_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			host_alias TEXT NOT NULL,
			host_name TEXT NOT NULL,
			port INTEGER NOT NULL DEFAULT 22,
			user TEXT NOT NULL DEFAULT '',
			last_connected DATETIME NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS favorite_hosts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			host_alias TEXT NOT NULL,
			host_name TEXT NOT NULL,
			port INTEGER NOT NULL DEFAULT 22,
			user TEXT NOT NULL DEFAULT '',
			tags TEXT NOT NULL DEFAULT '',
			group_name TEXT NOT NULL DEFAULT 'default'
		)`,
		`CREATE TABLE IF NOT EXISTS audit_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			action TEXT NOT NULL,
			detail TEXT NOT NULL,
			timestamp DATETIME NOT NULL
		)`,
	}

	for _, q := range queries {
		if _, err := d.db.Exec(q); err != nil {
			return err
		}
	}

	return nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) AddConnectionHistory(alias, hostname string, port int, user string) error {
	_, err := d.db.Exec(
		`INSERT INTO connection_history (host_alias, host_name, port, user, last_connected)
		 VALUES (?, ?, ?, ?, ?)`,
		alias, hostname, port, user, time.Now(),
	)
	return err
}

func (d *Database) GetConnectionHistory() ([]ConnectionHistory, error) {
	rows, err := d.db.Query(
		`SELECT id, host_alias, host_name, port, user, last_connected
		 FROM connection_history ORDER BY last_connected DESC LIMIT 200`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []ConnectionHistory
	for rows.Next() {
		var h ConnectionHistory
		if err := rows.Scan(&h.ID, &h.HostAlias, &h.HostName, &h.Port, &h.User, &h.LastConnected); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}

func (d *Database) DeleteConnectionHistory(id int64) error {
	_, err := d.db.Exec(`DELETE FROM connection_history WHERE id = ?`, id)
	return err
}

func (d *Database) SearchConnectionHistory(query string) ([]ConnectionHistory, error) {
	rows, err := d.db.Query(
		`SELECT id, host_alias, host_name, port, user, last_connected
		 FROM connection_history
		 WHERE host_alias LIKE ? OR host_name LIKE ? OR user LIKE ?
		 ORDER BY last_connected DESC LIMIT 50`,
		"%"+query+"%", "%"+query+"%", "%"+query+"%",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []ConnectionHistory
	for rows.Next() {
		var h ConnectionHistory
		if err := rows.Scan(&h.ID, &h.HostAlias, &h.HostName, &h.Port, &h.User, &h.LastConnected); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}

func (d *Database) AddFavoriteHost(alias, hostname string, port int, user, tags, group string) error {
	_, err := d.db.Exec(
		`INSERT INTO favorite_hosts (host_alias, host_name, port, user, tags, group_name)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		alias, hostname, port, user, tags, group,
	)
	return err
}

func (d *Database) GetFavoriteHosts() ([]FavoriteHost, error) {
	rows, err := d.db.Query(
		`SELECT id, host_alias, host_name, port, user, tags, group_name
		 FROM favorite_hosts ORDER BY id DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hosts []FavoriteHost
	for rows.Next() {
		var h FavoriteHost
		if err := rows.Scan(&h.ID, &h.HostAlias, &h.HostName, &h.Port, &h.User, &h.Tags, &h.Group); err != nil {
			return nil, err
		}
		hosts = append(hosts, h)
	}
	return hosts, nil
}

func (d *Database) UpdateFavoriteHost(id int64, alias, hostname string, port int, user, tags, group string) error {
	_, err := d.db.Exec(
		`UPDATE favorite_hosts SET host_alias=?, host_name=?, port=?, user=?, tags=?, group_name=?
		 WHERE id=?`,
		alias, hostname, port, user, tags, group, id,
	)
	return err
}

func (d *Database) DeleteFavoriteHost(id int64) error {
	_, err := d.db.Exec(`DELETE FROM favorite_hosts WHERE id = ?`, id)
	return err
}

func (d *Database) AddAuditLog(action, detail string) error {
	_, err := d.db.Exec(
		`INSERT INTO audit_logs (action, detail, timestamp) VALUES (?, ?, ?)`,
		action, detail, time.Now(),
	)
	return err
}

func (d *Database) GetAuditLogs() ([]AuditLog, error) {
	rows, err := d.db.Query(
		`SELECT id, action, detail, timestamp FROM audit_logs ORDER BY timestamp DESC LIMIT 500`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var l AuditLog
		if err := rows.Scan(&l.ID, &l.Action, &l.Detail, &l.Timestamp); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func (d *Database) GetGroups() ([]string, error) {
	rows, err := d.db.Query(`SELECT DISTINCT group_name FROM favorite_hosts ORDER BY group_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []string
	for rows.Next() {
		var g string
		if err := rows.Scan(&g); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}
