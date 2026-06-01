package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/example/iching-fiber-app/internal/domain"
)

type ReadingRepository struct{ db *sql.DB }
func NewReadingRepository(db *sql.DB) *ReadingRepository { return &ReadingRepository{db: db} }

func (r *ReadingRepository) Save(ctx context.Context, reading domain.Reading) (domain.Reading, error) {
	if err := ensureSchema(ctx, r.db); err != nil { return domain.Reading{}, err }
	linesJSON, _ := json.Marshal(reading.Lines)
	changingJSON, _ := json.Marshal(reading.ChangingLines)
	_, err := r.db.ExecContext(ctx, `INSERT INTO readings (id, question, method, primary_number, relating_number, changing_lines, lines, interpretation, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, reading.ID, reading.Question, string(reading.Method), reading.PrimaryNumber, reading.RelatingNumber, string(changingJSON), string(linesJSON), reading.Interpretation, reading.CreatedAt)
	return reading, err
}

func (r *ReadingRepository) List(ctx context.Context) ([]domain.Reading, error) {
	if err := ensureSchema(ctx, r.db); err != nil { return nil, err }
	rows, err := r.db.QueryContext(ctx, `SELECT id, question, method, primary_number, relating_number, changing_lines, lines, interpretation, created_at FROM readings ORDER BY created_at DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []domain.Reading
	for rows.Next() {
		item, err := scanReading(rows.Scan)
		if err != nil { return nil, err }
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *ReadingRepository) Get(ctx context.Context, id string) (domain.Reading, error) {
	if err := ensureSchema(ctx, r.db); err != nil { return domain.Reading{}, err }
	row := r.db.QueryRowContext(ctx, `SELECT id, question, method, primary_number, relating_number, changing_lines, lines, interpretation, created_at FROM readings WHERE id = ?`, id)
	return scanReading(row.Scan)
}

func ensureSchema(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS readings (
  id TEXT PRIMARY KEY,
  question TEXT NOT NULL,
  method TEXT NOT NULL,
  primary_number INTEGER NOT NULL,
  relating_number INTEGER NOT NULL,
  changing_lines TEXT NOT NULL,
  lines TEXT NOT NULL,
  interpretation TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL
)`)
	return err
}

type scanFn func(dest ...any) error
func scanReading(scan scanFn) (domain.Reading, error) {
	var item domain.Reading
	var method string
	var linesJSON, changingJSON string
	if err := scan(&item.ID, &item.Question, &method, &item.PrimaryNumber, &item.RelatingNumber, &changingJSON, &linesJSON, &item.Interpretation, &item.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) { return domain.Reading{}, errors.New("reading not found") }
		return domain.Reading{}, err
	}
	item.Method = domain.CastMethod(method)
	_ = json.Unmarshal([]byte(linesJSON), &item.Lines)
	_ = json.Unmarshal([]byte(changingJSON), &item.ChangingLines)
	return item, nil
}
