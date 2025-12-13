-- +goose Up
-- +goose StatementBegin

-- =============================================
-- SEED PROVIDERS
-- =============================================

INSERT INTO agc.providers (id, name, description, endpoint_url) VALUES
    ('00000000-0001-0000-0000-000000000001', 'OpenAI', 'OpenAI API - GPT models and more', 'https://api.openai.com/v1'),
    ('00000000-0002-0000-0000-000000000001', 'Anthropic', 'Anthropic API - Claude models', 'https://api.anthropic.com/v1');

-- =============================================
-- SEED MODEL SCHEMAS
-- =============================================

-- OpenAI Chat Options Schema
INSERT INTO agc.model_schemas (id, type, name, schema) VALUES
    ('00000000-0001-0001-0000-000000000001', 'options', 'openai_chat_options', '{
        "type": "object",
        "properties": {
            "temperature": {
                "type": "number",
                "minimum": 0,
                "maximum": 2,
                "default": 1,
                "description": "Sampling temperature between 0 and 2"
            },
            "max_tokens": {
                "type": "integer",
                "minimum": 1,
                "description": "Maximum number of tokens to generate"
            },
            "top_p": {
                "type": "number",
                "minimum": 0,
                "maximum": 1,
                "default": 1,
                "description": "Nucleus sampling probability"
            },
            "frequency_penalty": {
                "type": "number",
                "minimum": -2,
                "maximum": 2,
                "default": 0,
                "description": "Frequency penalty for token repetition"
            },
            "presence_penalty": {
                "type": "number",
                "minimum": -2,
                "maximum": 2,
                "default": 0,
                "description": "Presence penalty for new topics"
            },
            "stop": {
                "type": "array",
                "items": {"type": "string"},
                "maxItems": 4,
                "description": "Stop sequences"
            }
        }
    }');

-- OpenAI Chat Response Schema
INSERT INTO agc.model_schemas (id, type, name, schema) VALUES
    ('00000000-0001-0002-0000-000000000001', 'response', 'openai_chat_response', '{
        "type": "object",
        "properties": {
            "id": {"type": "string"},
            "object": {"type": "string"},
            "created": {"type": "integer"},
            "model": {"type": "string"},
            "choices": {
                "type": "array",
                "items": {
                    "type": "object",
                    "properties": {
                        "index": {"type": "integer"},
                        "message": {
                            "type": "object",
                            "properties": {
                                "role": {"type": "string"},
                                "content": {"type": "string"}
                            }
                        },
                        "finish_reason": {"type": "string"}
                    }
                }
            },
            "usage": {
                "type": "object",
                "properties": {
                    "prompt_tokens": {"type": "integer"},
                    "completion_tokens": {"type": "integer"},
                    "total_tokens": {"type": "integer"}
                }
            }
        }
    }');

-- Anthropic Chat Options Schema
INSERT INTO agc.model_schemas (id, type, name, schema) VALUES
    ('00000000-0002-0001-0000-000000000001', 'options', 'anthropic_chat_options', '{
        "type": "object",
        "properties": {
            "max_tokens": {
                "type": "integer",
                "minimum": 1,
                "description": "Maximum number of tokens to generate (required)"
            },
            "temperature": {
                "type": "number",
                "minimum": 0,
                "maximum": 1,
                "default": 1,
                "description": "Sampling temperature between 0 and 1"
            },
            "top_p": {
                "type": "number",
                "minimum": 0,
                "maximum": 1,
                "description": "Nucleus sampling probability"
            },
            "top_k": {
                "type": "integer",
                "minimum": 0,
                "description": "Top-k sampling"
            },
            "stop_sequences": {
                "type": "array",
                "items": {"type": "string"},
                "description": "Custom stop sequences"
            }
        },
        "required": ["max_tokens"]
    }');

-- Anthropic Chat Response Schema
INSERT INTO agc.model_schemas (id, type, name, schema) VALUES
    ('00000000-0002-0002-0000-000000000001', 'response', 'anthropic_chat_response', '{
        "type": "object",
        "properties": {
            "id": {"type": "string"},
            "type": {"type": "string"},
            "role": {"type": "string"},
            "model": {"type": "string"},
            "content": {
                "type": "array",
                "items": {
                    "type": "object",
                    "properties": {
                        "type": {"type": "string"},
                        "text": {"type": "string"}
                    }
                }
            },
            "stop_reason": {"type": "string"},
            "stop_sequence": {"type": "string"},
            "usage": {
                "type": "object",
                "properties": {
                    "input_tokens": {"type": "integer"},
                    "output_tokens": {"type": "integer"}
                }
            }
        }
    }');

-- =============================================
-- SEED OPENAI MODELS
-- =============================================

-- GPT-4o Series
INSERT INTO agc.models (id, model_key, name, description, provider_id, options_schema_id, response_schema_id, request_url) VALUES
    ('00000001-0001-0001-0000-000000000001', 'gpt-4o', 'GPT-4o', 'Most advanced multimodal model. Faster and cheaper than GPT-4 Turbo with stronger vision capabilities.', '00000000-0001-0000-0000-000000000001', '00000000-0001-0001-0000-000000000001', '00000000-0001-0002-0000-000000000001', 'https://api.openai.com/v1/chat/completions'),
    ('00000001-0001-0002-0000-000000000001', 'gpt-4o-mini', 'GPT-4o Mini', 'Affordable small model for fast, lightweight tasks. Cheaper and more capable than GPT-3.5 Turbo.', '00000000-0001-0000-0000-000000000001', '00000000-0001-0001-0000-000000000001', '00000000-0001-0002-0000-000000000001', 'https://api.openai.com/v1/chat/completions'),
    ('00000001-0001-0003-0000-000000000001', 'chatgpt-4o-latest', 'ChatGPT-4o Latest', 'Dynamic model continuously updated to the current version of GPT-4o in ChatGPT.', '00000000-0001-0000-0000-000000000001', '00000000-0001-0001-0000-000000000001', '00000000-0001-0002-0000-000000000001', 'https://api.openai.com/v1/chat/completions');

-- GPT-4 Turbo Series
INSERT INTO agc.models (id, model_key, name, description, provider_id, options_schema_id, response_schema_id, request_url) VALUES
    ('00000001-0002-0001-0000-000000000001', 'gpt-4-turbo', 'GPT-4 Turbo', 'GPT-4 Turbo with Vision. 128k context window, training data up to Dec 2023.', '00000000-0001-0000-0000-000000000001', '00000000-0001-0001-0000-000000000001', '00000000-0001-0002-0000-000000000001', 'https://api.openai.com/v1/chat/completions'),
    ('00000001-0002-0002-0000-000000000001', 'gpt-4-turbo-preview', 'GPT-4 Turbo Preview', 'GPT-4 Turbo preview model for testing latest features.', '00000000-0001-0000-0000-000000000001', '00000000-0001-0001-0000-000000000001', '00000000-0001-0002-0000-000000000001', 'https://api.openai.com/v1/chat/completions');

-- GPT-4 Base Series
INSERT INTO agc.models (id, model_key, name, description, provider_id, options_schema_id, response_schema_id, request_url) VALUES
    ('00000001-0003-0001-0000-000000000001', 'gpt-4', 'GPT-4', 'Original GPT-4 model with 8k context window.', '00000000-0001-0000-0000-000000000001', '00000000-0001-0001-0000-000000000001', '00000000-0001-0002-0000-000000000001', 'https://api.openai.com/v1/chat/completions'),
    ('00000001-0003-0002-0000-000000000001', 'gpt-4-32k', 'GPT-4 32k', 'GPT-4 with extended 32k context window.', '00000000-0001-0000-0000-000000000001', '00000000-0001-0001-0000-000000000001', '00000000-0001-0002-0000-000000000001', 'https://api.openai.com/v1/chat/completions');

-- O1 Reasoning Models
INSERT INTO agc.models (id, model_key, name, description, provider_id, options_schema_id, response_schema_id, request_url) VALUES
    ('00000001-0004-0001-0000-000000000001', 'o1', 'O1', 'Reasoning model designed to solve hard problems across domains. High intelligence for complex tasks.', '00000000-0001-0000-0000-000000000001', '00000000-0001-0001-0000-000000000001', '00000000-0001-0002-0000-000000000001', 'https://api.openai.com/v1/chat/completions'),
    ('00000001-0004-0002-0000-000000000001', 'o1-mini', 'O1 Mini', 'Faster and cheaper reasoning model, particularly good at coding, math, and science.', '00000000-0001-0000-0000-000000000001', '00000000-0001-0001-0000-000000000001', '00000000-0001-0002-0000-000000000001', 'https://api.openai.com/v1/chat/completions'),
    ('00000001-0004-0003-0000-000000000001', 'o1-preview', 'O1 Preview', 'Preview version of O1 reasoning model for early access to new capabilities.', '00000000-0001-0000-0000-000000000001', '00000000-0001-0001-0000-000000000001', '00000000-0001-0002-0000-000000000001', 'https://api.openai.com/v1/chat/completions');

-- O3 Reasoning Models
INSERT INTO agc.models (id, model_key, name, description, provider_id, options_schema_id, response_schema_id, request_url) VALUES
    ('00000001-0005-0001-0000-000000000001', 'o3-mini', 'O3 Mini', 'Latest reasoning model with improved performance. Fast and cost-effective for reasoning tasks.', '00000000-0001-0000-0000-000000000001', '00000000-0001-0001-0000-000000000001', '00000000-0001-0002-0000-000000000001', 'https://api.openai.com/v1/chat/completions');

-- GPT-3.5 Turbo Series
INSERT INTO agc.models (id, model_key, name, description, provider_id, options_schema_id, response_schema_id, request_url) VALUES
    ('00000001-0006-0001-0000-000000000001', 'gpt-3.5-turbo', 'GPT-3.5 Turbo', 'Fast, inexpensive model for simple tasks. 16k context window.', '00000000-0001-0000-0000-000000000001', '00000000-0001-0001-0000-000000000001', '00000000-0001-0002-0000-000000000001', 'https://api.openai.com/v1/chat/completions'),
    ('00000001-0006-0002-0000-000000000001', 'gpt-3.5-turbo-16k', 'GPT-3.5 Turbo 16k', 'GPT-3.5 Turbo with extended 16k context window.', '00000000-0001-0000-0000-000000000001', '00000000-0001-0001-0000-000000000001', '00000000-0001-0002-0000-000000000001', 'https://api.openai.com/v1/chat/completions');

-- =============================================
-- SEED ANTHROPIC MODELS
-- =============================================

-- Claude 4.5 Series (Latest)
INSERT INTO agc.models (id, model_key, name, description, provider_id, options_schema_id, response_schema_id, request_url) VALUES
    ('00000002-0001-0001-0000-000000000001', 'claude-sonnet-4-5-20250929', 'Claude Sonnet 4.5', 'Smart model for complex agents and coding. Best balance of intelligence, speed, and cost.', '00000000-0002-0000-0000-000000000001', '00000000-0002-0001-0000-000000000001', '00000000-0002-0002-0000-000000000001', 'https://api.anthropic.com/v1/messages'),
    ('00000002-0001-0002-0000-000000000001', 'claude-haiku-4-5-20251001', 'Claude Haiku 4.5', 'Fastest model with near-frontier intelligence.', '00000000-0002-0000-0000-000000000001', '00000000-0002-0001-0000-000000000001', '00000000-0002-0002-0000-000000000001', 'https://api.anthropic.com/v1/messages'),
    ('00000002-0001-0003-0000-000000000001', 'claude-opus-4-5-20251101', 'Claude Opus 4.5', 'Premium model combining maximum intelligence with practical performance.', '00000000-0002-0000-0000-000000000001', '00000000-0002-0001-0000-000000000001', '00000000-0002-0002-0000-000000000001', 'https://api.anthropic.com/v1/messages');

-- Claude 4 Series
INSERT INTO agc.models (id, model_key, name, description, provider_id, options_schema_id, response_schema_id, request_url) VALUES
    ('00000002-0002-0001-0000-000000000001', 'claude-opus-4-20250514', 'Claude Opus 4', 'Most powerful Claude 4 model for highly complex tasks.', '00000000-0002-0000-0000-000000000001', '00000000-0002-0001-0000-000000000001', '00000000-0002-0002-0000-000000000001', 'https://api.anthropic.com/v1/messages'),
    ('00000002-0002-0002-0000-000000000001', 'claude-sonnet-4-20250514', 'Claude Sonnet 4', 'High intelligence model balancing capability with speed.', '00000000-0002-0000-0000-000000000001', '00000000-0002-0001-0000-000000000001', '00000000-0002-0002-0000-000000000001', 'https://api.anthropic.com/v1/messages'),
    ('00000002-0002-0003-0000-000000000001', 'claude-opus-4-1-20250805', 'Claude Opus 4.1', 'Enhanced Opus 4 with improved capabilities.', '00000000-0002-0000-0000-000000000001', '00000000-0002-0001-0000-000000000001', '00000000-0002-0002-0000-000000000001', 'https://api.anthropic.com/v1/messages');

-- Claude 3.7 Series
INSERT INTO agc.models (id, model_key, name, description, provider_id, options_schema_id, response_schema_id, request_url) VALUES
    ('00000002-0003-0001-0000-000000000001', 'claude-3-7-sonnet-20250219', 'Claude 3.7 Sonnet', 'Claude 3.7 Sonnet with extended thinking support.', '00000000-0002-0000-0000-000000000001', '00000000-0002-0001-0000-000000000001', '00000000-0002-0002-0000-000000000001', 'https://api.anthropic.com/v1/messages');

-- Claude 3.5 Series (Legacy)
INSERT INTO agc.models (id, model_key, name, description, provider_id, options_schema_id, response_schema_id, request_url) VALUES
    ('00000002-0004-0001-0000-000000000001', 'claude-3-5-haiku-20241022', 'Claude 3.5 Haiku', 'Fast and compact model for near-instant responsiveness.', '00000000-0002-0000-0000-000000000001', '00000000-0002-0001-0000-000000000001', '00000000-0002-0002-0000-000000000001', 'https://api.anthropic.com/v1/messages');

-- Claude 3 Series (Legacy)
INSERT INTO agc.models (id, model_key, name, description, provider_id, options_schema_id, response_schema_id, request_url) VALUES
    ('00000002-0005-0001-0000-000000000001', 'claude-3-haiku-20240307', 'Claude 3 Haiku', 'Fastest and most compact legacy model.', '00000000-0002-0000-0000-000000000001', '00000000-0002-0001-0000-000000000001', '00000000-0002-0002-0000-000000000001', 'https://api.anthropic.com/v1/messages');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Delete models first (foreign key constraints)
DELETE FROM agc.models WHERE provider_id IN (
    '00000000-0001-0000-0000-000000000001',
    '00000000-0002-0000-0000-000000000001'
);

-- Delete model schemas
DELETE FROM agc.model_schemas WHERE id IN (
    '00000000-0001-0001-0000-000000000001',
    '00000000-0001-0002-0000-000000000001',
    '00000000-0002-0001-0000-000000000001',
    '00000000-0002-0002-0000-000000000001'
);

-- Delete providers
DELETE FROM agc.providers WHERE id IN (
    '00000000-0001-0000-0000-000000000001',
    '00000000-0002-0000-0000-000000000001'
);

-- +goose StatementEnd
