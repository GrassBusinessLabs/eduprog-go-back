CREATE TABLE IF NOT EXISTS public.eduprogresults(
    id serial PRIMARY KEY,
    "type" varchar(10),
    code integer,
    definition varchar(500),
    eduprog_id integer REFERENCES public.eduprog(id) ON DELETE CASCADE
);