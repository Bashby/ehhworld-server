-- Connect to database
\connect ehhworld


-- user
CREATE TABLE player_account (
	-- base columns
	id SERIAL PRIMARY KEY,
	uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
	created timestamptz NOT NULL DEFAULT now(),
	created_by text,
	last_modified timestamptz,
	last_modified_by text,
	deleted timestamptz,
	deleted_by text,

	-- table columns
	name TEXT NOT NULL
);