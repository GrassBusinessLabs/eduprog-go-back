CREATE TABLE IF NOT EXISTS discipline(
    id serial PRIMARY KEY,
    "name" varchar(50),
    eduprog_id integer REFERENCES eduprog(id),
    created_date TIMESTAMP DEFAULT now(),
    updated_date TIMESTAMP DEFAULT now()
);