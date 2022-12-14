package database

import (
	"contacts_cqrs/models"
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type RepositoryPostgres struct {
	db *sql.DB
}

func CNDatabase(url string) (*RepositoryPostgres, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &RepositoryPostgres{db: db}, err
}

func (r *RepositoryPostgres) Close() {
	r.db.Close()
}

func (r *RepositoryPostgres) InsertCts(ctx context.Context, ct *models.Contact) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO contacts(name, lastname, image, email, phone, status, updateData) VALUES($1, $2, $3, $4, $5, $6, $7)", ct.Name, ct.Lastname, ct.Image, ct.Email, ct.Phone, ct.Status, ct.UpdateData)
	return err
}

func (r *RepositoryPostgres) ListContact(ctx context.Context) ([]*models.Contact, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, lastname, image, email, phone, status, createdAt, updateData FROM contacts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cts := []*models.Contact{}
	for rows.Next() {
		ct := &models.Contact{}
		if err := rows.Scan(&ct.ID, &ct.Name, &ct.Lastname, &ct.Image, &ct.Email, &ct.Phone, &ct.Status, &ct.CreatedAt, &ct.UpdateData); err != nil {
			return nil, err
		}
		cts = append(cts, ct)
	}

	return cts, nil
}

func (r *RepositoryPostgres) UpdateCts(ctx context.Context, id string, ct *models.Contact) error {
	_, err := r.db.ExecContext(ctx, "UPDATE contacts SET name=$1, lastname=$2, image=$3, email=$4, phone=$5, status=$6, updateData=$7 WHERE id= $8", ct.Name, ct.Lastname, ct.Image, ct.Email, ct.Phone, ct.Status, ct.UpdateData, id)
	return err
}

func (r *RepositoryPostgres) DeleteCts(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM contacts WHERE id = $1", id)
	return err
}
