package service

import (
	"ahv-worldwide/internal/db"
)

func GetSettings() (map[string]string, error) {
	rows, err := db.DB.Query(`SELECT key, value FROM site_settings`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settings := map[string]string{}
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		settings[k] = v
	}
	return settings, nil
}

func UpdateSettings(settings map[string]string) error {
	for k, v := range settings {
		_, err := db.DB.Exec(`
			INSERT INTO site_settings (key, value) VALUES ($1, $2)
			ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value`, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
