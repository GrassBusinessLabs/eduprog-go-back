CREATE TABLE IF NOT EXISTS eduprog
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    "name"       TEXT,
    education_level        TEXT,
    stage     TEXT,
    speciality     TEXT,
    knowledge_field     TEXT,
    user_id     INTEGER,
    created_date DATETIME,
    updated_date DATETIME,
    deleted_date DATETIME NULL
);
