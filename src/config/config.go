package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var Strdbconn = ""
var Apiport = 0
var MaxConn = 0

// LoadEnvironment - Function that load the environment variables
func LoadEnvironment() {
	var erro error

	if erro = godotenv.Load(); erro != nil {
		log.Fatal(erro)
	}

	Apiport, erro = strconv.Atoi(os.Getenv("APIPORT"))
	if erro != nil {
		Apiport = 8080
	}

	MaxConn, erro = strconv.Atoi(os.Getenv("DBPOOLMAXCONNS"))
	if erro != nil {
		MaxConn = 20
	}

	Strdbconn = fmt.Sprintf(`
		user=` + os.Getenv("DBUSER") + `
		password=` + os.Getenv("DBPASSWORD") + `
		host=` + os.Getenv("DBHOST") + `
		port=` + os.Getenv("DBPORT") + `
		dbname=` + os.Getenv("DBNAME") + `
		sslmode=disable`)

}
