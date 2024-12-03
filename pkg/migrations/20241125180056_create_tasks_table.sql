-- +goose Up
CREATE SEQUENCE IF NOT EXISTS tasks_id_seq;

CREATE TABLE IF NOT EXISTS public.tasks
(
    id integer NOT NULL DEFAULT nextval('tasks_id_seq'::regclass),
    title character varying(255) NOT NULL,
    description text NOT NULL,
    due_date date,
    priority character varying(50),
    file_path character varying(255),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_by integer,
    CONSTRAINT tasks_pkey PRIMARY KEY (id),
    CONSTRAINT tasks_created_by_fkey FOREIGN KEY (created_by)
        REFERENCES public.users (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)
    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.tasks
    OWNER TO roo;


-- +goose Down
DROP TABLE IF EXISTS tasks;
