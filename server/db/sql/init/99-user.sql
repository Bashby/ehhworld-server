-- Connect to database
\connect ehhworld


-- Create user
CREATE USER ehhworld_db_docker_user WITH PASSWORD 'ehhworlddb';

-- Grant user control of database, tables, indexes
GRANT ALL PRIVILEGES ON DATABASE ehhworld TO ehhworld_db_docker_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ehhworld_db_docker_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ehhworld_db_docker_user;