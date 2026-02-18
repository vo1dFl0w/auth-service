CREATE TABLE users (
    user_id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    provider_id TEXT,
    oauth_provider TEXT,

    CONSTRAINT check_auth_method CHECK (
        (password_hash IS NOT NULL) OR 
        (provider_id IS NOT NULL AND oauth_provider IS NOT NULL)
    )
);