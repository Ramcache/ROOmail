-- +goose Up
CREATE TABLE IF NOT EXISTS public.tasks_users
(
    task_id integer NOT NULL,
    user_id integer NOT NULL,
    assigned_at timestamp DEFAULT CURRENT_TIMESTAMP,
    sent_by integer,
    CONSTRAINT tasks_users_pkey PRIMARY KEY (task_id, user_id),
    CONSTRAINT tasks_users_sent_by_fkey FOREIGN KEY (sent_by)
        REFERENCES public.users (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT tasks_users_task_id_fkey FOREIGN KEY (task_id)
        REFERENCES public.tasks (id)
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT tasks_users_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES public.users (id)
        ON UPDATE NO ACTION
        ON DELETE CASCADE
)
    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.tasks_users
    OWNER TO roo;

-- +goose Down
DROP TABLE IF EXISTS task_recipients;
