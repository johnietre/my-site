CREATE TABLE IF NOT EXISTS products (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  webpage TEXT NOT NULL,
  app_store_link TEXT NOT NULL,
  play_store_link TEXT NOT NULL,
  hidden BOOLEAN NOT NULL
);

create TABLE IF NOT EXISTS issues (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  product_id INTEGER,
  email TEXT NOT NULL,
  reason TEXT NOT NULL,
  subject TEXT NOT NULL,
  description TEXT NOT NULL,
  -- Second precision
  created_at INT64 NOT NULL,
  started_at INT64 NOT NULL,
  resolved_at INT64 NOT NULL,
  ip TEXT NOT NULL,
  FOREIGN KEY(product_id) REFERENCES apps(id)
);
