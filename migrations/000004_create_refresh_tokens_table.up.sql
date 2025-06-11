CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY,                                -- ID токена
    token TEXT NOT NULL,                                -- сам refresh token (строка/JWT)
    user_id UUID REFERENCES users(id) ON DELETE CASCADE, -- связанный пользователь
    expires_at TIMESTAMP NOT NULL,                      -- когда истекает
    created_at TIMESTAMP DEFAULT now(),                 -- когда создан
    revoked BOOLEAN DEFAULT FALSE                       -- отозван или нет
);