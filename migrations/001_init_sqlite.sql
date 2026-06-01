CREATE TABLE IF NOT EXISTS readings (
  id TEXT PRIMARY KEY,
  question TEXT NOT NULL,
  method TEXT NOT NULL,
  primary_number INTEGER NOT NULL,
  relating_number INTEGER NOT NULL,
  changing_lines TEXT NOT NULL,
  lines TEXT NOT NULL,
  interpretation TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS hexagrams (
  number INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  judgment TEXT NOT NULL,
  image TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS hexagram_lines (
  hexagram_number INTEGER NOT NULL,
  line_position INTEGER NOT NULL,
  line_text TEXT NOT NULL,
  PRIMARY KEY (hexagram_number, line_position)
);
