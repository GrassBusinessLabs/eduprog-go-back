CREATE TABLE IF NOT EXISTS public.competencies_base (
    id SERIAL PRIMARY KEY,
    "type" varchar(50) NOT NULL,
    definition varchar(500),
    specialty varchar(50)
);
