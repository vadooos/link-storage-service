package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"link-storage-service/model"

	_ "github.com/lib/pq"
)

type LinkRepository struct {
	db *sql.DB
}

func Open(dsn string) (*LinkRepository, error) {
	db, err := sql.Open("postgres", dsn)
	return &LinkRepository{db: db}, err
}

func (r *LinkRepository) Close() error {
	return r.db.Close()
}

func (r *LinkRepository) List(limit int, offset int) ([]model.Link, error) {
	const q = `
		SELECT id, short_code, original_url, created_at::text, visits 
		FROM links 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(q, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list links: %w", err)
	}
	defer rows.Close()
	var links []model.Link
	for rows.Next() {
		var link model.Link
		err := rows.Scan(&link.ID, &link.ShortCode, &link.OriginalURL, &link.CreatedAt, &link.Visits)
		if err != nil {
			return nil, fmt.Errorf("scan links: %w", err)
		}
		links = append(links, link)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}
	return links, nil
}

func (r *LinkRepository) Insert(link model.Link) (model.Link, error) {
	const q = `
		INSERT INTO links (short_code, original_url) 
		VALUES ($1, $2) 
		RETURNING id, short_code, original_url, created_at::text, visits
	`
	err := r.db.QueryRow(q, link.ShortCode, link.OriginalURL).Scan(&link.ID, &link.ShortCode, &link.OriginalURL, &link.CreatedAt, &link.Visits)
	if err != nil {
		return model.Link{}, fmt.Errorf("create link: %w", err)
	}
	return link, err
}

func (r *LinkRepository) GetByShortCode(shortCode string) (model.Link, error) {
	const q = `SELECT id, short_code, original_url, created_at::text, visits FROM links WHERE short_code = $1`
	var link model.Link
	err := r.db.QueryRow(q, shortCode).Scan(&link.ID, &link.ShortCode, &link.OriginalURL, &link.CreatedAt, &link.Visits)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Link{}, ErrNotFound
		}
		return model.Link{}, fmt.Errorf("get link: %w", err)
	}
	return link, nil
}

func (r *LinkRepository) IncrementVisits(shortCode string) (int64, error) {
	const q = `UPDATE links SET visits = visits + 1 WHERE short_code = $1 RETURNING visits`
	var visits int64
	err := r.db.QueryRow(q, shortCode).Scan(&visits)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, fmt.Errorf("increment visits: %w", err)
	}
	return visits, nil
}

func (r *LinkRepository) Delete(shortCode string) error {
	const q = `DELETE FROM links WHERE short_code = $1`
	res, err := r.db.Exec(q, shortCode)
	if err != nil {
		return fmt.Errorf("delete link: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete link rows: %w", err)
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}
