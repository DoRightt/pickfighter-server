\c pickfighter;

-- Create the users table
CREATE TABLE IF NOT EXISTS auth.users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    rank VARCHAR(60),
    claim VARCHAR(60) DEFAULT 'USER',
    flags INTEGER NOT NULL DEFAULT 0,
    created_at INTEGER NOT NULL DEFAULT (EXTRACT(epoch FROM now()))::INTEGER,
    updated_at INTEGER
);

-- Create a unique index for the 'name' field in the users table
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_name ON auth.users(name);

-- Create the user_credentials table
CREATE TABLE IF NOT EXISTS auth.user_credentials (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    email VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL,
    salt VARCHAR NOT NULL,
    token VARCHAR UNIQUE,
    token_type VARCHAR,
    active BOOLEAN DEFAULT false,
    token_expire INTEGER,
    created_at INTEGER NOT NULL DEFAULT (EXTRACT(epoch FROM now()))::INTEGER,
    updated_at INTEGER,
    CONSTRAINT fk_user_credentials_user_id FOREIGN KEY (user_id) 
        REFERENCES auth.users(user_id) ON DELETE CASCADE
);

-- Create a unique index for the 'email' field in the user_credentials table
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_credentials_email ON auth.user_credentials(email);

-- Create an index for the 'user_id' field in the user_credentials table
CREATE INDEX IF NOT EXISTS idx_user_credentials_user_id ON auth.user_credentials(user_id);