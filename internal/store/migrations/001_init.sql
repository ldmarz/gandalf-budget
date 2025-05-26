CREATE TABLE months (
  id INTEGER PRIMARY KEY,
  year INT NOT NULL,
  month INT NOT NULL,
  finalized BOOLEAN DEFAULT 0
);
CREATE TABLE categories (
  id INTEGER PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,
  color TEXT NOT NULL          -- Tailwind colour class
);
CREATE TABLE budget_lines (
  id INTEGER PRIMARY KEY,
  month_id INT NOT NULL REFERENCES months(id) ON DELETE CASCADE,
  category_id INT NOT NULL REFERENCES categories(id),
  label TEXT NOT NULL,
  expected REAL NOT NULL
);
CREATE TABLE actual_lines (
  id INTEGER PRIMARY KEY,
  budget_line_id INT NOT NULL REFERENCES budget_lines(id) ON DELETE CASCADE,
  actual REAL NOT NULL DEFAULT 0
);
CREATE TABLE annual_snaps (
  id INTEGER PRIMARY KEY,
  month_id INT NOT NULL REFERENCES months(id) UNIQUE,
  snap_json TEXT NOT NULL,
  created_at DATETIME NOT NULL
);
