package dblayer

var tables []string = []string{`
CREATE TABLE IF NOT EXISTS strazh (
	id          SERIAL PRIMARY KEY,
	uid 		INT ,
	origin        TEXT,
	short       TEXT,
	deleted		BOOLEAN
)`,
	`CREATE TABLE IF NOT EXISTS users (
		id          SERIAL PRIMARY KEY,
			name 		TEXT ,
		password        TEXT
)`}
