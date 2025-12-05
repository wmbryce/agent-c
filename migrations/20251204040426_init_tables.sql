-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- Add UUID extension

-- Create schema
CREATE SCHEMA IF NOT EXISTS agc;

-- Create tables
-- Create users table
CREATE TABLE IF NOT EXISTS agc.providers (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name VARCHAR (255) NOT NULL,
    description TEXT NOT NULL,
    endpoint_url VARCHAR (255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS agc.model_schemas (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    type VARCHAR (255) NOT NULL CHECK (type IN ('options', 'response')),
    name VARCHAR (255) NOT NULL,
    schema JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL
);


CREATE TABLE IF NOT EXISTS agc.models (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    model_key VARCHAR (255) NOT NULL UNIQUE,
    name VARCHAR (255) NOT NULL,
    description TEXT NOT NULL,
    provider_id UUID NOT NULL REFERENCES agc.providers (id) ON DELETE CASCADE,
    options_schema_id UUID NOT NULL REFERENCES agc.model_schemas (id) ON DELETE CASCADE,
    response_schema_id UUID NOT NULL REFERENCES agc.model_schemas (id) ON DELETE CASCADE,
    request_url VARCHAR (255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL
);

CREATE TABLE agc.sellers (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    wallet_address VARCHAR (255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL
);

CREATE TABLE agc.consumers (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    wallet_address VARCHAR (255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL
);

CREATE TABLE agc.api_keys (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    api_key VARCHAR (255) NOT NULL UNIQUE,
    tokens_available INT NOT NULL,
    provider_id UUID NOT NULL REFERENCES agc.providers (id) ON DELETE CASCADE,
    seller_id UUID NOT NULL REFERENCES agc.sellers (id) ON DELETE CASCADE,
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
DROP TABLE IF EXISTS agc.model_schemas;
DROP TABLE IF EXISTS agc.providers;
DROP TABLE IF EXISTS agc.sellers;
DROP TABLE IF EXISTS agc.consumers;

DROP SCHEMA IF EXISTS agc;
-- +goose StatementEnd
