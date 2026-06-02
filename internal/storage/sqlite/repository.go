package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
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
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS readings (
		  id TEXT PRIMARY KEY,
		  question TEXT NOT NULL,
		  method TEXT NOT NULL,
		  language TEXT,
		  lines_json TEXT NOT NULL,
		  primary_number INTEGER NOT NULL,
		  relating_number INTEGER NOT NULL,
		  changing_lines_json TEXT NOT NULL,
		  interpretation_json TEXT NOT NULL,
		  reflection_rating INTEGER,
		  reflection_note TEXT,
		  reflection_created_at TEXT,
		  created_at TEXT NOT NULL
		);`,
		`ALTER TABLE readings ADD COLUMN language TEXT;`,
		`ALTER TABLE readings ADD COLUMN reflection_rating INTEGER;`,
		`ALTER TABLE readings ADD COLUMN reflection_note TEXT;`,
		`ALTER TABLE readings ADD COLUMN reflection_created_at TEXT;`,
		`UPDATE readings SET language = 'cs' WHERE language IS NULL OR TRIM(language) = '';`,
	}
	for _, stmt := range stmts {
		if _, err := r.db.Exec(stmt); err != nil {
			msg := strings.ToLower(err.Error())
			if strings.Contains(msg, "duplicate column name") {
				continue
			}
		}
	}
}

func (r *ReadingRepository) Save(ctx context.Context, in domain.Reading) (domain.Reading, error) {
	linesJSON, _ := json.Marshal(in.Lines)
	changingJSON, _ := json.Marshal(in.ChangingLines)
	interpJSON, _ := json.Marshal(in.Interpretation)

	lang := strings.TrimSpace(in.Language)
	if lang == "" {
		lang = "cs"
	}

	_, err := r.db.ExecContext(ctx, `
INSERT INTO readings (
  id, question, method, language, lines_json, primary_number, relating_number,
  changing_lines_json, interpretation_json, reflection_rating, reflection_note,
  reflection_created_at, created_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		in.ID,
		in.Question,
		string(in.Method),
		lang,
		string(linesJSON),
		in.PrimaryNumber,
		in.RelatingNumber,
		string(changingJSON),
		string(interpJSON),
		nullableRating(in.Reflection),
		nullableNote(in.Reflection),
		nullableTime(in.Reflection.CreatedAt),
		in.CreatedAt.Format(time.RFC3339Nano),
	)
	if err != nil {
		return domain.Reading{}, fmt.Errorf("insert sqlite reading: %w", err)
	}
	in.Language = lang
	return in, nil
}

func (r *ReadingRepository) List(ctx context.Context) ([]domain.Reading, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT
  id,
  question,
  method,
  COALESCE(language, 'cs') AS language,
  lines_json,
  primary_number,
  relating_number,
  changing_lines_json,
  interpretation_json,
  reflection_rating,
  reflection_note,
  reflection_created_at,
  created_at
FROM readings
ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var out []domain.Reading
	for rows.Next() {
		item, err := scanReading(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *ReadingRepository) Get(ctx context.Context, id string) (domain.Reading, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT
  id,
  question,
  method,
  COALESCE(language, 'cs') AS language,
  lines_json,
  primary_number,
  relating_number,
  changing_lines_json,
  interpretation_json,
  reflection_rating,
  reflection_note,
  reflection_created_at,
  created_at
FROM readings
WHERE id = ?`, id)
	return scanReading(row.Scan)
}

func (r *ReadingRepository) SaveReflection(ctx context.Context, id string, rating int, note string, createdAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `
UPDATE readings
SET reflection_rating = ?, reflection_note = ?, reflection_created_at = ?
WHERE id = ?`,
		rating,
		note,
		createdAt.Format(time.RFC3339Nano),
		id,
	)
	return err
}

type scanner func(dest ...any) error

func scanReading(scan scanner) (domain.Reading, error) {
	var item domain.Reading
	var method, linesJSON, changingJSON, interpJSON, createdAt string
	var language sql.NullString
	var reflectionRating sql.NullInt64
	var reflectionNote sql.NullString
	var reflectionCreatedAt sql.NullString

	if err := scan(
		&item.ID,
		&item.Question,
		&method,
		&language,
		&linesJSON,
		&item.PrimaryNumber,
		&item.RelatingNumber,
		&changingJSON,
		&interpJSON,
		&reflectionRating,
		&reflectionNote,
		&reflectionCreatedAt,
		&createdAt,
	); err != nil {
		return domain.Reading{}, err
	}

	item.Method = domain.CastMethod(method)
	if language.Valid && strings.TrimSpace(language.String) != "" {
		item.Language = language.String
	} else {
		item.Language = "cs"
	}

	_ = json.Unmarshal([]byte(linesJSON), &item.Lines)
	_ = json.Unmarshal([]byte(changingJSON), &item.ChangingLines)
	_ = json.Unmarshal([]byte(interpJSON), &item.Interpretation)

	if t, err := time.Parse(time.RFC3339Nano, createdAt); err == nil {
		item.CreatedAt = t
	}
	if reflectionRating.Valid {
		item.Reflection.Rating = int(reflectionRating.Int64)
	}
	if reflectionNote.Valid {
		item.Reflection.Note = reflectionNote.String
	}
	if reflectionCreatedAt.Valid {
		if t, err := time.Parse(time.RFC3339Nano, reflectionCreatedAt.String); err == nil {
			item.Reflection.CreatedAt = &t
		}
	}

	return item, nil
}

func nullableTime(t *time.Time) any {
	if t == nil {
		return nil
	}
	return t.Format(time.RFC3339Nano)
}

func nullableRating(r domain.Reflection) any {
	if r.CreatedAt == nil && r.Rating == 0 && r.Note == "" {
		return nil
	}
	return r.Rating
}

func nullableNote(r domain.Reflection) any {
	if r.Note == "" {
		return nil
	}
	return r.Note
}
