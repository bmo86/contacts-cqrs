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
