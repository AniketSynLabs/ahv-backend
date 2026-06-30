package service

import (
	"time"

	"ahv-worldwide/internal/db"
	"ahv-worldwide/internal/model"
)

func CreateLead(l model.Lead) (model.Lead, error) {
	var id int
	err := db.DB.QueryRow(`
		INSERT INTO leads (name, email, phone, country, inquiry, message, source)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		l.Name, l.Email, l.Phone, l.Country, l.Inquiry, l.Message, l.Source,
	).Scan(&id)
	if err != nil {
		return model.Lead{}, err
	}
	l.ID = id
	l.Status = "New"
	l.CreatedAt = time.Now().Format(time.RFC3339)
	return l, nil
}

func ListLeads(status string) ([]model.Lead, error) {
	query := `
		SELECT id, name, email, phone, country, inquiry, message, source, status, created_at
		FROM leads`
	args := []any{}
	if status != "" {
		query += ` WHERE status = $1`
		args = append(args, status)
	}
	query += ` ORDER BY created_at DESC`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leads []model.Lead
	for rows.Next() {
		var l model.Lead
		var createdAt time.Time
		if err := rows.Scan(
			&l.ID, &l.Name, &l.Email, &l.Phone, &l.Country,
			&l.Inquiry, &l.Message, &l.Source, &l.Status, &createdAt,
		); err != nil {
			return nil, err
		}
		l.CreatedAt = createdAt.Format(time.RFC3339)
		leads = append(leads, l)
	}
	if leads == nil {
		leads = []model.Lead{}
	}
	return leads, nil
}

func UpdateLeadStatus(id, status string) error {
	_, err := db.DB.Exec(`UPDATE leads SET status = $1 WHERE id = $2`, status, id)
	return err
}

func DeleteLead(id string) error {
	_, err := db.DB.Exec(`DELETE FROM leads WHERE id = $1`, id)
	return err
}
