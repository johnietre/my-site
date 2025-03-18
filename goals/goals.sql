CREATE TABLE IF NOT EXISTS goals_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  what TEXT NOT NULL,
  type INTEGER NOT NULL,
  parent INTEGER,
  completed BOOLEAN,
  hidden BOOLEAN NOT NULL
);

-- type-None = 0
-- type-UListItem = 1
-- type-OListItem = 2
-- type-Checkbox = 3
-- type-Text = 4
-- type-Header1 = 5
-- type-Header2 = 6
-- type-Header3 = 7
-- type-Header4 = 8
-- type-Header5 = 9
-- type-Header6 = 10
