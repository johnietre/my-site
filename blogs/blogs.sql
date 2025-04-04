CREATE TABLE IF NOT EXISTS categories (
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS blogs (
  id INTEGER PRIMARY KEY,
  title TEXT NOT NULL,
  -- format |AUTHOR 1|AUTHOR 2|...|
  authors TEXT NOT NULL,
  -- format |CAT 1|CAT 2|...|
  categories TEXT NOT NULL,
  timestamp INTEGER NOT NULL,
  tz_offset INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS edits (
  blog_id INTEGER NOT NULL,
  description TEXT NOT NULL,
  timestamp INTEGER NOT NULL,
  tz_offset INTEGER NOT NULL,
  prev_hash TEXT NOT NULL,
  hash TEXT NOT NULL,
  FOREIGN KEY(blog_id) REFERENCES blogs(id)
);
