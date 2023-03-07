CREATE TABLE IF NOT EXISTS public.eduprogcompetencies(
    id SERIAL PRIMARY KEY,
    competency_id INTEGER NOT NULL REFERENCES public.competencies_base (id),
    eduprog_id INTEGER NOT NULL REFERENCES public.eduprog (id),
    "type" varchar(10),
    code INTEGER,
    redefinition VARCHAR(200)
);
