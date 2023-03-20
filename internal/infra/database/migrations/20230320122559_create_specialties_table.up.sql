CREATE TABLE IF NOT EXISTS public.specialties (
    code varchar(10) PRIMARY KEY,
    "name" varchar(100),
    kf_code varchar(10),
    knowledge_field varchar(100)
);
