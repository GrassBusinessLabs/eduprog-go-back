CREATE TABLE IF NOT EXISTS sessions(
    user_id integer NOT NULL,
    uuid varchar(50) NOT NULL,
    CONSTRAINT auths_pkey PRIMARY KEY (user_id,uuid)
    );