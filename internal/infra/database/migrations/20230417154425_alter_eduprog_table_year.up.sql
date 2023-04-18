ALTER TABLE public.eduprog
    ADD COLUMN approval_year INTEGER;
--     ADD COLUMN child_of INTEGER REFERENCES eduprog(id) ON DELETE CASCADE DEFAULT null;