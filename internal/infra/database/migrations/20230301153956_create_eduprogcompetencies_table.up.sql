CREATE TABLE IF NOT EXISTS public.eduprogcompetencies(
    id SERIAL PRIMARY KEY,
    competency_id INTEGER REFERENCES public.competencies_base (id) ON DELETE CASCADE,
    eduprog_id INTEGER NOT NULL REFERENCES public.eduprog (id) ON DELETE CASCADE,
    "type" varchar(10),
    code INTEGER,
    definition VARCHAR(500)
);

ALTER TABLE public.eduprogcompetencies
ALTER COLUMN competency_id DROP NOT NULL;
