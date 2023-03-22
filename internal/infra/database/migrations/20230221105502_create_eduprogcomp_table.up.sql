CREATE TABLE IF NOT EXISTS public.eduprogcomp(
    id serial PRIMARY KEY,
    code varchar(50) NOT NULL ,
    "name" varchar(50) NOT NULL ,
    credits float8 NOT NULL,
    control_type varchar(50),
    "type" varchar(50),
    block_num varchar(50),
    block_name varchar(100),
    "category" varchar(50),
    eduprog_id integer REFERENCES eduprog(id) ON DELETE CASCADE,
    created_date TIMESTAMP DEFAULT now(),
    updated_date TIMESTAMP DEFAULT now(),
    deleted_date TIMESTAMP NULL
    );