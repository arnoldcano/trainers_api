CREATE DATABASE trainers_api;
\c trainers_api;
CREATE TABLE trainers (
    id SERIAL PRIMARY KEY,
    email TEXT,
    phone TEXT,
    first_name TEXT,
    last_name TEXT
);
