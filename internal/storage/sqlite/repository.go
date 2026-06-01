package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/example/iching-fiber-app/internal/domain"
)

type ReadingRepository struct{ db *sql.DB }

func NewReadingRepository(db *sql.DB) *ReadingRepository {
	r := &ReadingRepository{db: db}
	r.mustInit()
	return r
}

func (r *ReadingRepository) mustInit() {
	const schema = `
CREATE TABLE IF NOT EXISTS readings (
  id TEXT PRIMARY KEY,
  question TEXT NOT NULL,
  method TEXT NOT NULL,
  lines_json TEXT NOT NULL,
  primary_number INTEGER NOT NULL,
  relating_number INTEGER NOT NULL,
  changing_lines_json TEXT NOT NULL,
  interpretation_json TEXT NOT NULL,
  created_at TEXT NOT NULL
);`
	if _, err := r.db.Exec(schema); err != nil { panic(err) }
}

func (r *ReadingRepository) Save(ctx context.Context, in domain.Reading) (domain.Reading, error) {
	linesJSON, _ := json.Marshal(in.Lines)
	changingJSON, _ := json.Marshal(in.ChangingLines)
	interpJSON, _ := json.Marshal(in.Interpretation)
	_, err := r.db.ExecContext(ctx, `
INSERT INTO readings (id, question, method, lines_json, primary_number, relating_number, changing_lines_json, interpretation_json, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		in.ID, in.Question, string(in.Method), string(linesJSON), in.PrimaryNumber, in.RelatingNumber, string(changingJSON), string(interpJSON), in.CreatedAt.Format(time.RFC3339Nano),
	)
	if err != nil { return domain.Reading{}, fmt.Errorf("insert sqlite reading: %w", err) }
	return in, nil
}

func (r *ReadingRepository) List(ctx context.Context) ([]domain.Reading, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, question, method, lines_json, primary_number, relating_number, changing_lines_json, interpretation_json, created_at FROM readings ORDER BY created_at DESC`)
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
	row := r.db.QueryRowContext(ctx, `SELECT id, question, method, lines_json, primary_number, relating_number, changing_lines_json, interpretation_json, created_at FROM readings WHERE id = ?`, id)
	return scanReading(row.Scan)
}

type scanner func(dest ...any) error

func scanReading(scan scanner) (domain.Reading, error) {
	var item domain.Reading
	var method, linesJSON, changingJSON, interpJSON, createdAt string
	if err := scan(&item.ID, &item.Question, &method, &linesJSON, &item.PrimaryNumber, &item.RelatingNumber, &changingJSON, &interpJSON, &createdAt); err != nil {
		return domain.Reading{}, err
	}
	item.Method = domain.CastMethod(method)
	_ = json.Unmarshal([]byte(linesJSON), &item.Lines)
	_ = json.Unmarshal([]byte(changingJSON), &item.ChangingLines)
	_ = json.Unmarshal([]byte(interpJSON), &item.Interpretation)
	if t, err := time.Parse(time.RFC3339Nano, createdAt); err == nil { item.CreatedAt = t }
	return item, nil
}
