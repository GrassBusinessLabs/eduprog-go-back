CREATE TABLE IF NOT EXISTS eduprogcomp
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    code      TEXT,
    "name"       TEXT,
    credits       INTEGER,
    control_type     TEXT,
    "type"     TEXT,
    sub_type     TEXT,
    category     TEXT,
    eduprog_id INTEGER,
    created_date DATETIME,
    updated_date DATETIME,
    deleted_date DATETIME NULL
);
