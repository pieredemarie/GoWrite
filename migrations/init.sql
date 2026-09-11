CREATE TABLE IF NOT EXISTS articles (
    id              SERIAL PRIMARY KEY,
    slug            TEXT UNIQUE NOT NULL,
    title           TEXT NOT NULL,
    contentmd       TEXT NOT NULL,
    contenthtml     TEXT NOT NULL,
    authortokenhash TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);