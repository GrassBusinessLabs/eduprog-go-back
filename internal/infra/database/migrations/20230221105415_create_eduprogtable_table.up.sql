CREATE TABLE IF NOT EXISTS public.eduprog(
    id serial PRIMARY KEY,
    "name" varchar(50) NOT NULL,
    education_level varchar(50),
    stage varchar(50),
    speciality_code varchar(10),
    speciality varchar(100),
    kf_code varchar(10),
    knowledge_field varchar(100),
    user_id integer REFERENCES users(id),
    created_date TIMESTAMP DEFAULT now(),
    updated_date TIMESTAMP DEFAULT now(),
    deleted_date TIMESTAMP NULL
    );