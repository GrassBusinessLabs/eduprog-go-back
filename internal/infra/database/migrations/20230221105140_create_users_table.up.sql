CREATE TABLE IF NOT EXISTS public.users(
    id SERIAL PRIMARY KEY,
    "name" VARCHAR(50),
    email VARCHAR(50),
    "password" VARCHAR(1500),
    created_date TIMESTAMP DEFAULT now(),
    updated_date TIMESTAMP DEFAULT now(),
    deleted_date TIMESTAMP NULL
    );