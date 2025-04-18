CREATE TABLE IF NOT EXISTS users(
    id uuid PRIMARY KEY,
    name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    username character varying(255) NOT NULL,
    role character varying(50) NOT NULL,
    password character varying(255) NOT NULL,
    status character varying(50) NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);

