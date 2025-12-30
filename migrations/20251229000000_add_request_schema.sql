-- +goose Up
-- +goose StatementBegin

-- =============================================
-- ADD REQUEST_SCHEMA COLUMN (replaces request_format)
-- =============================================

-- Add new request_schema column
ALTER TABLE agc.providers ADD COLUMN request_schema JSONB DEFAULT '{}';

-- Drop old request_format column if it exists
ALTER TABLE agc.providers DROP COLUMN IF EXISTS request_format;

-- =============================================
-- UPDATE OPENAI PROVIDER WITH REQUEST SCHEMA
-- =============================================

UPDATE agc.providers SET
  request_schema = '{
    "model_field": "model",
    "messages_field": "messages"
  }'
WHERE name = 'OpenAI';

-- =============================================
-- UPDATE ANTHROPIC PROVIDER WITH REQUEST SCHEMA
-- =============================================

UPDATE agc.providers SET
  request_schema = '{
    "model_field": "model",
    "messages_field": "messages"
  }'
WHERE name = 'Anthropic';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE agc.providers DROP COLUMN IF EXISTS request_schema;
ALTER TABLE agc.providers ADD COLUMN request_format VARCHAR(50);

-- +goose StatementEnd
