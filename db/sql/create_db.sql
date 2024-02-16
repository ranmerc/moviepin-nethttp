CREATE DATABASE moviepin;

\c moviepin;

-- Enable the UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE Users (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE Movies (
    movie_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    release_date DATE,
    genre TEXT,
    director TEXT,
    description TEXT
);

CREATE TABLE Reviews (
    review_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES Users(user_id),
    movie_id UUID REFERENCES Movies(movie_id),
    rating INTEGER,
    review_text TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Lists (
    list_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES Users(user_id),
    title TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ListItems (
    list_item_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    list_id UUID REFERENCES Lists(list_id),
    movie_id UUID REFERENCES Movies(movie_id),
    position INTEGER
);
