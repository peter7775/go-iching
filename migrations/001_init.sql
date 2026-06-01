CREATE TABLE IF NOT EXISTS readings (
  id TEXT PRIMARY KEY,
  question TEXT NOT NULL,
  method TEXT NOT NULL,
  primary_number INTEGER NOT NULL,
  relating_number INTEGER NOT NULL,
  changing_lines JSONB NOT NULL,
  lines JSONB NOT NULL,
  interpretation TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
