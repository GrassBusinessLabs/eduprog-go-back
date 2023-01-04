CREATE TABLE IF NOT EXISTS users
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    "name"       TEXT,
    email        TEXT,
    password     TEXT,
    created_date DATETIME,
    updated_date DATETIME,
    deleted_date DATETIME NULL
);
