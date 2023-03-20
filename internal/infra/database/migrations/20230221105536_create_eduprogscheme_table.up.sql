CREATE TABLE IF NOT EXISTS public.eduprogscheme(
    id serial PRIMARY KEY,
    semester_num integer NOT NULL,
    discipline varchar(50),
    discipline_id integer ,
    eduprog_id integer REFERENCES eduprog(id) ON DELETE CASCADE,
    eduprogcomp_id integer REFERENCES eduprogcomp(id) ON DELETE CASCADE,
    credits_per_semester float8,
    created_date TIMESTAMP DEFAULT now(),
    updated_date TIMESTAMP DEFAULT now()
    );