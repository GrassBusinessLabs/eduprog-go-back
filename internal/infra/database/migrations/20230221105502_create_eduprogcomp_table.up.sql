CREATE TABLE IF NOT EXISTS public.eduprogcomp(
    id serial PRIMARY KEY,
    code varchar(50) NOT NULL,
    "name" varchar(50) NOT NULL,
    credits integer NOT NULL,
    control_type varchar(50),
    "type" varchar(50),
    sub_type varchar(50),
    "category" varchar(50),
    eduprog_id integer REFERENCES eduprog(id),
    created_date TIMESTAMP DEFAULT now(),
    updated_date TIMESTAMP DEFAULT now(),
    deleted_date TIMESTAMP NULL
    );