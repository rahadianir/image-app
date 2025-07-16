CREATE SCHEMA IF NOT EXISTS "project";

CREATE TABLE IF NOT EXISTS project.images (
    id UUID PRIMARY KEY,
    file_name TEXT NOT NULL,
    url TEXT NOT NULL,
    file_size BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);