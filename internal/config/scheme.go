package config

// Scheme represents the application configuration scheme.
type Scheme struct {
	Env  string
	DB   *DB
	Http *Http
}

// DB is service Data base connection params
type DB struct {
	Driver string
	Addr   string
}

type Http struct {
	TCP string
}
