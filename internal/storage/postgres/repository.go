package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/example/iching-fiber-app/internal/domain"
)

type ReadingRepository struct {
	db *sql.DB
}

func NewReadingRepository(db *sql.DB) *ReadingRepository {
	return &ReadingRepository{db: db}
}

func (r *ReadingRepository) Save(ctx context.Context, in domain.Reading) (domain.Reading, error) {
	linesJSON, _ := json.Marshal(in.Lines)
	changingJSON, _ := json.Marshal(in.ChangingLines)
	interpJSON, _ := json.Marshal(in.Interpretation)

	_, err := r.db.ExecContext(ctx, `
INSERT INTO readings (
	id, question, method, language, lines_json, primary_number, relating_number,
	changing_lines_json, interpretation_json,
	reflection_rating, reflection_note, reflection_created_at, created_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		in.ID, in.Question, string(in.Method), in.Language,
		string(linesJSON), in.PrimaryNumber, in.RelatingNumber,
		string(changingJSON), string(interpJSON),
		nullableRating(in.Reflection), nullableNote(in.Reflection), nullableTime(in.Reflection.CreatedAt),
		in.CreatedAt,
	)
	if err != nil {
		return domain.Reading{}, fmt.Errorf("insert postgres reading: %w", err)
	}
	return in, nil
}

func (r *ReadingRepository) List(ctx context.Context) ([]domain.Reading, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT id, question, method, language, lines_json, primary_number, relating_number,
       changing_lines_json, interpretation_json,
       reflection_rating, reflection_note, reflection_created_at, created_at
FROM readings
ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
SELECT id, question, method, language, lines_json, primary_number, relating_number,
       changing_lines_json, interpretation_json,
       reflection_rating, reflection_note, reflection_created_at, created_at
FROM readings
WHERE id = $1`, id)
	return scanReading(row.Scan)
}

func (r *ReadingRepository) SaveReflection(ctx context.Context, id string, rating int, note string, createdAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `
UPDATE readings
SET reflection_rating = $1, reflection_note = $2, reflection_created_at = $3
WHERE id = $4`, rating, note, createdAt, id)
	return err
}

type scanner func(dest ...any) error

func scanReading(scan scanner) (domain.Reading, error) {
	var item domain.Reading
	var method, language, linesJSON, changingJSON, interpJSON string
	var reflectionRating sql.NullInt64
	var reflectionNote sql.NullString
	var reflectionCreatedAt sql.NullTime
	var createdAt time.Time

	if err := scan(
		&item.ID, &item.Question, &method, &language, &linesJSON,
		&item.PrimaryNumber, &item.RelatingNumber,
		&changingJSON, &interpJSON,
		&reflectionRating, &reflectionNote, &reflectionCreatedAt, &createdAt,
	); err != nil {
		return domain.Reading{}, err
	}

	item.Method = domain.CastMethod(method)
	item.Language = language
	_ = json.Unmarshal([]byte(linesJSON), &item.Lines)
	_ = json.Unmarshal([]byte(changingJSON), &item.ChangingLines)
	_ = json.Unmarshal([]byte(interpJSON), &item.Interpretation)
	item.CreatedAt = createdAt

	if reflectionRating.Valid {
		item.Reflection.Rating = int(reflectionRating.Int64)
	}
	if reflectionNote.Valid {
		item.Reflection.Note = reflectionNote.String
	}
	if reflectionCreatedAt.Valid {
		t := reflectionCreatedAt.Time
		item.Reflection.CreatedAt = &t
	}

	return item, nil
}

func nullableTime(t *time.Time) any {
	if t == nil {
		return nil
	}
	return *t
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
