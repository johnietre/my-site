CREATE TABLE categories (
  name TEXT NOT NULL UNIQUE,
);

CREATE TABLE blogs (
  id INTEGER PRIMARY KEY AUTO INCREMENT,
  title TEXT NOT NULL,
  timestamp INTEGER NOT NULL,
  tz_offset INTEGER NOT NULL,
  -- Format: |cat1|cat2|...|
  categories TEXT NOT NULL
);

CREATE TABLE edits (
  blog_id INTEGER NOT NULL,
  timestamp INTEGER NOT NULL,
  tz_offset INTEGER NOT NULL,
  FOREIGN KEY(blog_id) REFERENCES blogs(id)
);
