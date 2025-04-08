CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create a superuser with admin privileges
-- Password is 'admin123' (hashed with bcrypt)
INSERT INTO users (username, password_hash, email, role, created_at, updated_at)
VALUES (
    'admin',
    '$2a$12$RHDj3N3c/Q8wtaZmbCeI.uMV2YOJA4aVJyOkZEZDpZmyl63kf1y7a', -- bcrypt hash for 'admin123'
    'admin@example.com',
    'admin',
    NOW(),
    NOW()
);