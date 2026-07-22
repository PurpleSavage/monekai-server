-- 000001_create_tables.up.sql

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    credits INT DEFAULT 0,
    created_at TIMESTAMP,
    photo_url TEXT,
    CONSTRAINT users_external_id_key UNIQUE (external_id),
    CONSTRAINT users_email_key UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    refresh_token VARCHAR NOT NULL,
    user_agent VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT sessions_refresh_token_key UNIQUE (refresh_token),
    CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);

CREATE TABLE IF NOT EXISTS samples (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    sample_name VARCHAR(100) NOT NULL,
    prompt TEXT,
    initial_audio_url TEXT,
    duration INT,
    model_version VARCHAR(30) NOT NULL,
    output_format VARCHAR(10) NOT NULL,
    prediction_id VARCHAR,
    status VARCHAR(20) DEFAULT 'processing',
    created_at TIMESTAMP,
    CONSTRAINT samples_prediction_id_key UNIQUE (prediction_id),
    CONSTRAINT samples_model_version_check CHECK (model_version IN ('stereo-melody-large', 'stereo-large', 'melody-large', 'large')),
    CONSTRAINT samples_output_format_check CHECK (output_format IN ('mp3', 'wav')),
    CONSTRAINT samples_status_check CHECK (status IN ('starting', 'processing', 'succeeded', 'failed', 'canceled')),
    CONSTRAINT fk_samples_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_samples_user_id ON samples(user_id);

CREATE TABLE IF NOT EXISTS sample_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sample_id UUID NOT NULL,
    effects JSONB,
    final_audio_url TEXT,
    created_at TIMESTAMP,
    CONSTRAINT fk_sample_versions_sample FOREIGN KEY (sample_id) REFERENCES samples(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_sample_versions_sample_id ON sample_versions(sample_id);

CREATE TABLE IF NOT EXISTS shared_samples (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sample_id UUID,
    sample_version_id UUID,
    user_id UUID NOT NULL,
    likes INT DEFAULT 0,
    downloads INT DEFAULT 0,
    created_at TIMESTAMP,
    CONSTRAINT fk_shared_samples_sample FOREIGN KEY (sample_id) REFERENCES samples(id) ON DELETE SET NULL,
    CONSTRAINT fk_shared_samples_version FOREIGN KEY (sample_version_id) REFERENCES sample_versions(id) ON DELETE SET NULL,
    CONSTRAINT fk_shared_samples_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_shared_samples_sample_id ON shared_samples(sample_id);
CREATE INDEX IF NOT EXISTS idx_shared_samples_version_id ON shared_samples(sample_version_id);
CREATE INDEX IF NOT EXISTS idx_shared_samples_user_id ON shared_samples(user_id);

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    type VARCHAR(50),
    title VARCHAR(255),
    message TEXT,
    status VARCHAR(20) DEFAULT 'unread',
    reference_id VARCHAR(255),
    created_at TIMESTAMP,
    CONSTRAINT fk_notifications_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);


CREATE TABLE IF NOT EXISTS payments (

    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL,

    credit_package_id UUID NOT NULL,

    paddle_transaction_id VARCHAR NOT NULL UNIQUE,

    paddle_price_id VARCHAR NOT NULL,

    credits_purchased INT NOT NULL,

    amount NUMERIC(10,2) NOT NULL,

    currency VARCHAR(3) NOT NULL,

    status VARCHAR(30) NOT NULL,

    raw_payload JSONB,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_payment_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_payment_package
        FOREIGN KEY(credit_package_id)
        REFERENCES credit_packages(id)
);
CREATE INDEX IF NOT EXISTS idx_payments_user_id ON payments(user_id);

CREATE TABLE IF NOT EXISTS credit_packages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    paddle_price_id VARCHAR NOT NULL UNIQUE,

    name VARCHAR(100) NOT NULL,

    credits INT NOT NULL,

    price NUMERIC(10,2) NOT NULL,

    currency VARCHAR(3) NOT NULL DEFAULT 'USD',

    active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);