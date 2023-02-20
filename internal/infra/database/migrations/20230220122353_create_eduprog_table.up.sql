CREATE TABLE IF NOT EXISTS eduprog(
    id serial PRIMARY KEY,
    "name" varchar(50) NOT NULL,
    education_level varchar(50),
    stage varchar(50),
    speciality varchar(50),
    knowledge_field varchar(50),
    user_id integer REFERENCES users(id),
    created_date TIMESTAMP DEFAULT now(),
    updated_date TIMESTAMP DEFAULT now(),
    deleted_date TIMESTAMP NULL
);