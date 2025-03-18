CREATE TABLE IF NOT EXISTS categories (
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS blogs (
  id INTEGER PRIMARY KEY,
  title TEXT NOT NULL,
  timestamp INTEGER NOT NULL,
  tz_offset INTEGER NOT NULL,
  -- Format: |cat1|cat2|...|
  categories TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS edits (
  blog_id INTEGER NOT NULL,
  timestamp INTEGER NOT NULL,
  tz_offset INTEGER NOT NULL,
  FOREIGN KEY(blog_id) REFERENCES blogs(id)
);
