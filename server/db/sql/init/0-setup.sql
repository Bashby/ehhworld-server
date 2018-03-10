-- Create database
CREATE DATABASE ehhworld;

-- Connect to database
\connect ehhworld

-- Load database extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; -- uuid generation
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- trigram fuzzy text searching