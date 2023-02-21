CREATE TABLE IF NOT EXISTS public.eduprogscheme(
    id serial PRIMARY KEY,
    semester_num integer NOT NULL,
    discipline varchar(50),
    eduprog_id integer REFERENCES eduprog(id),
    eduprogcomp_id integer REFERENCES eduprogcomp(id),
    credits_per_semester integer,
    created_date TIMESTAMP DEFAULT now(),
    updated_date TIMESTAMP DEFAULT now()
    );