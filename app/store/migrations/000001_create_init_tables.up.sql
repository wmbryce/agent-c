-- Add UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE models (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    model_key VARCHAR (255) NOT NULL UNIQUE,
    options_schema JSONB NOT NULL,
    response_schema JSONB NOT NULL,
    request_url VARCHAR (255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL,
);

CREATE TABLE api_keys (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    api_key VARCHAR (255) NOT NULL UNIQUE,
    tokens_available INT NOT NULL,
    model_id UUID NOT NULL REFERENCES models (id) ON DELETE CASCADE,
    owner_id UUID NOT NULL REFERENCES owners (id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL,
)    

CREATE TABLE owners (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL,
)

CREATE TABLE consumers (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL,
)

-- Add indexes