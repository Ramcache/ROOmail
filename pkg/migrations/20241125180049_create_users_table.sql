-- +goose Up
CREATE TABLE IF NOT EXISTS public.users
(
    id integer NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    username character varying(50) COLLATE pg_catalog."default" NOT NULL,
    password_hash text COLLATE pg_catalog."default" NOT NULL,
    role character varying(50) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_username_key UNIQUE (username)
)

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.users
    OWNER to roo;
-- +goose Down
DROP TABLE IF EXISTS users;
