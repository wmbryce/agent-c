-- +goose Up
-- +goose StatementBegin

-- =============================================
-- ADD PROVIDER CONFIGURATION COLUMNS
-- =============================================

ALTER TABLE agc.providers ADD COLUMN auth_type VARCHAR(50);
ALTER TABLE agc.providers ADD COLUMN auth_header VARCHAR(100);
ALTER TABLE agc.providers ADD COLUMN extra_headers JSONB DEFAULT '{}';
ALTER TABLE agc.providers ADD COLUMN request_defaults JSONB DEFAULT '{}';
ALTER TABLE agc.providers ADD COLUMN response_mapping JSONB DEFAULT '{}';

-- =============================================
-- UPDATE OPENAI PROVIDER CONFIG
-- =============================================

UPDATE agc.providers SET
  auth_type = 'bearer',
  auth_header = 'Authorization',
  extra_headers = '{}',
  request_defaults = '{}',
  response_mapping = '{
    "id": "$.id",
    "model": "$.model",
    "content": "$.choices[0].message.content",
    "role": "$.choices[0].message.role",
    "finish_reason": "$.choices[0].finish_reason",
    "prompt_tokens": "$.usage.prompt_tokens",
    "completion_tokens": "$.usage.completion_tokens",
    "total_tokens": "$.usage.total_tokens"
  }'
WHERE name = 'OpenAI';

-- =============================================
-- UPDATE ANTHROPIC PROVIDER CONFIG
-- =============================================

UPDATE agc.providers SET
  auth_type = 'api_key',
  auth_header = 'x-api-key',
  extra_headers = '{"anthropic-version": "2023-06-01"}',
  request_defaults = '{"max_tokens": 1024}',
  response_mapping = '{
    "id": "$.id",
    "model": "$.model",
    "content": "$.content[0].text",
    "role": "$.role",
    "finish_reason": "$.stop_reason",
    "prompt_tokens": "$.usage.input_tokens",
    "completion_tokens": "$.usage.output_tokens",
    "total_tokens": null
  }'
WHERE name = 'Anthropic';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE agc.providers DROP COLUMN IF EXISTS auth_type;
ALTER TABLE agc.providers DROP COLUMN IF EXISTS auth_header;
ALTER TABLE agc.providers DROP COLUMN IF EXISTS extra_headers;
ALTER TABLE agc.providers DROP COLUMN IF EXISTS request_defaults;
ALTER TABLE agc.providers DROP COLUMN IF EXISTS response_mapping;

-- +goose StatementEnd
