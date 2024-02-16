CREATE TABLE IF NOT EXISTS Movies (
    movie_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    release_date DATE,
    genre TEXT, -- make enum
    director TEXT,
    description TEXT
);

CREATE TABLE IF NOT EXISTS Reviews (
    review_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES Users(user_id),
    movie_id UUID REFERENCES Movies(movie_id),
    rating DOUBLE PRECISION,
    review_text TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);