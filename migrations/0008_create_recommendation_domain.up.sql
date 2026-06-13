CREATE TABLE vehicle_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE vehicles (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL REFERENCES vehicle_categories(id),
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    version VARCHAR(100),
    year INTEGER NOT NULL CHECK (year BETWEEN 1980 AND 2050),
    fuel_type VARCHAR(50) NOT NULL,
    transmission VARCHAR(50) NOT NULL,
    price_min NUMERIC(12,2) NOT NULL,
    price_max NUMERIC(12,2) NOT NULL,
    seats INTEGER,
    trunk_capacity INTEGER,
    consumption_city NUMERIC(5,2),
    consumption_highway NUMERIC(5,2),
    description TEXT,
    strengths TEXT,
    weaknesses TEXT,
    match_profile JSONB NOT NULL DEFAULT '{}'::jsonb,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_vehicles_active ON vehicles(active);
CREATE INDEX idx_vehicles_category ON vehicles(category_id);
CREATE INDEX idx_vehicles_brand_model ON vehicles(brand, model);

CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    question_text TEXT NOT NULL,
    question_type VARCHAR(50) NOT NULL DEFAULT 'SINGLE_CHOICE',
    weight NUMERIC(4,3) NOT NULL DEFAULT 1,
    display_order INTEGER NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_questions_active_order ON questions(active, display_order);

CREATE TABLE answer_options (
    id SERIAL PRIMARY KEY,
    question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    option_text VARCHAR(500) NOT NULL,
    score_profile JSONB NOT NULL,
    display_order INTEGER NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (question_id, display_order)
);

CREATE TABLE recommendations (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    answers_snapshot JSONB NOT NULL,
    summary TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_recommendations_user_created ON recommendations(user_id, created_at DESC);

CREATE TABLE user_answers (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recommendation_id BIGINT NOT NULL REFERENCES recommendations(id) ON DELETE CASCADE,
    question_id INTEGER NOT NULL REFERENCES questions(id),
    answer_option_id INTEGER NOT NULL REFERENCES answer_options(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (recommendation_id, question_id)
);

CREATE INDEX idx_user_answers_user_created ON user_answers(user_id, created_at DESC);

CREATE TABLE recommendation_items (
    id BIGSERIAL PRIMARY KEY,
    recommendation_id BIGINT NOT NULL REFERENCES recommendations(id) ON DELETE CASCADE,
    vehicle_id INTEGER NOT NULL REFERENCES vehicles(id),
    rank INTEGER NOT NULL,
    score NUMERIC(5,2) NOT NULL,
    reason TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (recommendation_id, rank),
    UNIQUE (recommendation_id, vehicle_id)
);

CREATE INDEX idx_recommendation_items_recommendation ON recommendation_items(recommendation_id, rank);
