package utils

import "fmt"

type postgresEndpoint struct {
	host, port, dbName, user, password string
}

func defaultFromEnv() *postgresEndpoint {
	return &postgresEndpoint{
		host:     "127.0.0.1",
		port:     "5432",
		dbName:   "postgres",
		user:     "postgres",
		password: "123456",
	}
}

func GetDbInfo() string {
	defaultDb := defaultFromEnv()
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		defaultDb.host,
		defaultDb.port,
		defaultDb.dbName,
		defaultDb.user,
		defaultDb.password,
	)
}
