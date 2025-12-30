-- +goose Up
-- +goose StatementBegin

-- =============================================
-- SEED GOOGLE GEMINI PROVIDER
-- =============================================

INSERT INTO agc.providers (
    id, name, description, endpoint_url,
    auth_type, auth_header, extra_headers,
    request_defaults, response_mapping, request_schema
) VALUES (
    '00000000-0003-0000-0000-000000000001',
    'Google Gemini',
    'Google Gemini API - Advanced multimodal models',
    'https://generativelanguage.googleapis.com/v1beta',
    'api_key',
    'x-goog-api-key',
    '{}',
    '{}',
    '{
        "id": "$.candidates[0].content.parts[0].text",
        "model": "$.modelVersion",
        "content": "$.candidates[0].content.parts[0].text",
        "role": "$.candidates[0].content.role",
        "finish_reason": "$.candidates[0].finishReason",
        "prompt_tokens": "$.usageMetadata.promptTokenCount",
        "completion_tokens": "$.usageMetadata.candidatesTokenCount",
        "total_tokens": "$.usageMetadata.totalTokenCount"
    }',
    '{
        "model_field": "",
        "messages_field": "contents",
        "message_transform": {
            "role_field": "role",
            "content_path": "parts[].text",
            "role_map": {
                "assistant": "model"
            }
        },
        "options_wrapper": "generationConfig",
        "options_rename": {
            "max_tokens": "maxOutputTokens",
            "top_p": "topP",
            "top_k": "topK",
            "stop": "stopSequences",
            "stop_sequences": "stopSequences"
        }
    }'
);

-- =============================================
-- SEED GEMINI MODEL SCHEMAS
-- =============================================

-- Gemini Chat Options Schema
INSERT INTO agc.model_schemas (id, type, name, schema) VALUES
    ('00000000-0003-0001-0000-000000000001', 'options', 'gemini_chat_options', '{
        "type": "object",
        "properties": {
            "temperature": {
                "type": "number",
                "minimum": 0,
                "maximum": 2,
                "default": 1,
                "description": "Sampling temperature"
            },
            "max_tokens": {
                "type": "integer",
                "minimum": 1,
                "description": "Maximum number of output tokens (maps to maxOutputTokens)"
            },
            "top_p": {
                "type": "number",
                "minimum": 0,
                "maximum": 1,
                "description": "Nucleus sampling probability (maps to topP)"
            },
            "top_k": {
                "type": "integer",
                "minimum": 1,
                "description": "Top-k sampling (maps to topK)"
            },
            "stop_sequences": {
                "type": "array",
                "items": {"type": "string"},
                "description": "Stop sequences (maps to stopSequences)"
            },
            "candidate_count": {
                "type": "integer",
                "minimum": 1,
                "maximum": 8,
                "description": "Number of response candidates to generate"
            }
        }
    }');

-- Gemini Chat Response Schema
INSERT INTO agc.model_schemas (id, type, name, schema) VALUES
    ('00000000-0003-0002-0000-000000000001', 'response', 'gemini_chat_response', '{
        "type": "object",
        "properties": {
            "candidates": {
                "type": "array",
                "items": {
                    "type": "object",
                    "properties": {
                        "content": {
                            "type": "object",
                            "properties": {
                                "role": {"type": "string"},
                                "parts": {
                                    "type": "array",
                                    "items": {
                                        "type": "object",
                                        "properties": {
                                            "text": {"type": "string"}
                                        }
                                    }
                                }
                            }
                        },
                        "finishReason": {"type": "string"},
                        "safetyRatings": {"type": "array"}
                    }
                }
            },
            "usageMetadata": {
                "type": "object",
                "properties": {
                    "promptTokenCount": {"type": "integer"},
                    "candidatesTokenCount": {"type": "integer"},
                    "totalTokenCount": {"type": "integer"}
                }
            },
            "modelVersion": {"type": "string"}
        }
    }');

-- =============================================
-- SEED GEMINI 3 MODELS (Latest)
-- =============================================

INSERT INTO agc.models (id, model_key, name, description, provider_id, options_schema_id, response_schema_id, request_url) VALUES
    ('00000003-0001-0001-0000-000000000001', 'gemini-3-pro-preview', 'Gemini 3 Pro Preview', 'The best model in the world for multimodal understanding, and the most powerful agentic and coding model.', '00000000-0003-0000-0000-000000000001', '00000000-0003-0001-0000-000000000001', '00000000-0003-0002-0000-000000000001', 'https://generativelanguage.googleapis.com/v1beta/models/gemini-3-pro-preview:generateContent'),
    ('00000003-0001-0002-0000-000000000001', 'gemini-3-flash-preview', 'Gemini 3 Flash Preview', 'Balanced model built for speed, scale, and frontier intelligence.', '00000000-0003-0000-0000-000000000001', '00000000-0003-0001-0000-000000000001', '00000000-0003-0002-0000-000000000001', 'https://generativelanguage.googleapis.com/v1beta/models/gemini-3-flash-preview:generateContent');

-- =============================================
-- SEED GEMINI 2.5 MODELS
-- =============================================

INSERT INTO agc.models (id, model_key, name, description, provider_id, options_schema_id, response_schema_id, request_url) VALUES
    ('00000003-0002-0001-0000-000000000001', 'gemini-2.5-pro', 'Gemini 2.5 Pro', 'Advanced reasoning model for complex problems in code, math, STEM, and document analysis with long context support.', '00000000-0003-0000-0000-000000000001', '00000000-0003-0001-0000-000000000001', '00000000-0003-0002-0000-000000000001', 'https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-pro:generateContent'),
    ('00000003-0002-0002-0000-000000000001', 'gemini-2.5-flash', 'Gemini 2.5 Flash', 'Best model in terms of price-performance, offering well-rounded capabilities for large-scale, low-latency tasks.', '00000000-0003-0000-0000-000000000001', '00000000-0003-0001-0000-000000000001', '00000000-0003-0002-0000-000000000001', 'https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent'),
    ('00000003-0002-0003-0000-000000000001', 'gemini-2.5-flash-lite', 'Gemini 2.5 Flash Lite', 'Fastest flash model optimized for cost-efficiency and high throughput.', '00000000-0003-0000-0000-000000000001', '00000000-0003-0001-0000-000000000001', '00000000-0003-0002-0000-000000000001', 'https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash-lite:generateContent'),
    ('00000003-0002-0004-0000-000000000001', 'gemini-2.5-pro-preview-tts', 'Gemini 2.5 Pro TTS', 'Gemini 2.5 Pro with text-to-speech capabilities.', '00000000-0003-0000-0000-000000000001', '00000000-0003-0001-0000-000000000001', '00000000-0003-0002-0000-000000000001', 'https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-pro-preview-tts:generateContent'),
    ('00000003-0002-0005-0000-000000000001', 'gemini-2.5-flash-preview-tts', 'Gemini 2.5 Flash TTS', 'Gemini 2.5 Flash with text-to-speech capabilities.', '00000000-0003-0000-0000-000000000001', '00000000-0003-0001-0000-000000000001', '00000000-0003-0002-0000-000000000001', 'https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash-preview-tts:generateContent');

-- =============================================
-- SEED GEMINI 2.0 MODELS
-- =============================================

INSERT INTO agc.models (id, model_key, name, description, provider_id, options_schema_id, response_schema_id, request_url) VALUES
    ('00000003-0003-0001-0000-000000000001', 'gemini-2.0-flash', 'Gemini 2.0 Flash', 'Second-generation workhorse model with 1M token context window.', '00000000-0003-0000-0000-000000000001', '00000000-0003-0001-0000-000000000001', '00000000-0003-0002-0000-000000000001', 'https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent'),
    ('00000003-0003-0002-0000-000000000001', 'gemini-2.0-flash-lite', 'Gemini 2.0 Flash Lite', 'Cost-efficient variant of Gemini 2.0 Flash.', '00000000-0003-0000-0000-000000000001', '00000000-0003-0001-0000-000000000001', '00000000-0003-0002-0000-000000000001', 'https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash-lite:generateContent');

-- =============================================
-- SEED GEMINI 1.5 MODELS (Legacy but still available)
-- =============================================

INSERT INTO agc.models (id, model_key, name, description, provider_id, options_schema_id, response_schema_id, request_url) VALUES
    ('00000003-0004-0001-0000-000000000001', 'gemini-1.5-pro', 'Gemini 1.5 Pro', 'Previous generation pro model with strong reasoning capabilities.', '00000000-0003-0000-0000-000000000001', '00000000-0003-0001-0000-000000000001', '00000000-0003-0002-0000-000000000001', 'https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-pro:generateContent'),
    ('00000003-0004-0002-0000-000000000001', 'gemini-1.5-flash', 'Gemini 1.5 Flash', 'Previous generation flash model optimized for speed.', '00000000-0003-0000-0000-000000000001', '00000000-0003-0001-0000-000000000001', '00000000-0003-0002-0000-000000000001', 'https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent'),
    ('00000003-0004-0003-0000-000000000001', 'gemini-1.5-flash-8b', 'Gemini 1.5 Flash 8B', 'Lightweight 8B parameter variant of Gemini 1.5 Flash.', '00000000-0003-0000-0000-000000000001', '00000000-0003-0001-0000-000000000001', '00000000-0003-0002-0000-000000000001', 'https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-8b:generateContent');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Delete Gemini models first (foreign key constraints)
DELETE FROM agc.models WHERE provider_id = '00000000-0003-0000-0000-000000000001';

-- Delete Gemini model schemas
DELETE FROM agc.model_schemas WHERE id IN (
    '00000000-0003-0001-0000-000000000001',
    '00000000-0003-0002-0000-000000000001'
);

-- Delete Gemini provider
DELETE FROM agc.providers WHERE id = '00000000-0003-0000-0000-000000000001';

-- +goose StatementEnd
