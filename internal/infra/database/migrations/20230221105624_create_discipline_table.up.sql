
CREATE TABLE IF NOT EXISTS public.discipline(
    id serial PRIMARY KEY,
    "name" varchar(50),
    eduprog_id integer REFERENCES eduprog(id),
    created_date TIMESTAMP DEFAULT now(),
    updated_date TIMESTAMP DEFAULT now()
    );
ALTER TABLE public.eduprogscheme ADD CONSTRAINT fk_disc_id FOREIGN KEY (discipline_id)
    references discipline(id) ON DELETE CASCADE;