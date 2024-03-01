CREATE TABLE apps (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  webpage TEXT NOT NULL,
  on_app_store BOOLEAN NOT NULL,
  on_play_store BOOLEAN NOT NULL
);

create TABLE issues (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  app_id INTEGER,
  email TEXT NOT NULL,
  reason TEXT NOT NULL,
  subject TEXT NOT NULL,
  description TEXT NOT NULL,
  replied_to BOOLEAN NOT NULL,
  ip TEXT NOT NULL,
  -- Second precision
  timestamp INTEGER NOT NULL,
  FOREIGN KEY(app_id) REFERENCES apps(id)
);
