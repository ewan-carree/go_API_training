package pkg

import "os"

func GetEnvVariables() (string, string, string, string, string, string) {
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	port := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("NAME")
	return dialect, host, port, user, password, dbname
}