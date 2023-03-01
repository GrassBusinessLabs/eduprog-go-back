CREATE TABLE IF NOT EXISTS public.eduprogcompetencies(
    id SERIAL PRIMARY KEY,
    competency_id INTEGER NOT NULL UNIQUE REFERENCES public.competencies_base (id),
    eduprog_id INTEGER NOT NULL REFERENCES public.eduprog (id),
    code INTEGER NOT NULL UNIQUE,
    redefinition VARCHAR(200)
);
