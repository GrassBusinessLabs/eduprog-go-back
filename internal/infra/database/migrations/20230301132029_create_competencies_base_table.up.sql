CREATE TABLE IF NOT EXISTS public.competencies_base (
    id SERIAL PRIMARY KEY,
    "type" varchar(50) NOT NULL,
    code integer,
    definition varchar(500),
    specialty varchar(50)
);
