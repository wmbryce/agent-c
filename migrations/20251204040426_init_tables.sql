-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- Add UUID extension

-- Create schema
CREATE SCHEMA IF NOT EXISTS agc;

-- Create tables
-- Create users table
CREATE TABLE IF NOT EXISTS agc.models (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    model_key VARCHAR (255) NOT NULL UNIQUE,
    options_schema JSONB NOT NULL,
    response_schema JSONB NOT NULL,
    request_url VARCHAR (255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL
);

CREATE TABLE agc.owners (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL
);

CREATE TABLE agc.consumers (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL
);

CREATE TABLE agc.api_keys (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    api_key VARCHAR (255) NOT NULL UNIQUE,
    tokens_available INT NOT NULL,
    model_id UUID NOT NULL REFERENCES agc.models (id) ON DELETE CASCADE,
    owner_id UUID NOT NULL REFERENCES agc.owners (id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- Delete tables
DROP TABLE IF EXISTS agc.api_keys;
DROP TABLE IF EXISTS agc.models;
DROP TABLE IF EXISTS agc.owners;
DROP TABLE IF EXISTS agc.consumers;

DROP SCHEMA IF EXISTS agc;
-- +goose StatementEnd
