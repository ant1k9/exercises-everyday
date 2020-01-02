package db

const (
	// InitialMigration is the script to create all the tables in the database
	InitialMigration = `
CREATE TABLE IF NOT EXISTS exercises (
	id		SERIAL PRIMARY KEY NOT NULL,
	type	VARCHAR(256) NOT NULL,
	repeats INTEGER NOT NULL,
	date	TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
	id		 SERIAL PRIMARY KEY NOT NULL,
	login	 VARCHAR(255) NOT NULL,
	password VARCHAR(128) NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
	user_id INTEGER NOT NULL,
	value VARCHAR(255) NOT NULL,
	FOREIGN KEY ( user_id ) REFERENCES users( id ) ON DELETE CASCADE
);
`
)
