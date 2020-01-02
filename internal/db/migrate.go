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
`
)
