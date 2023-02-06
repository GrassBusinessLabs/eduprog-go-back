CREATE TABLE IF NOT EXISTS eduprogscheme
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    semester_num INTEGER,
    discipline TEXT,
    eduprog_id INTEGER,
    eduprogcomp_id INTEGER,
    credits_per_semester INTEGER,
    created_date DATETIME,
    updated_date DATETIME
);
