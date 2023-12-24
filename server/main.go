package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"server/db"
	"server/routes"
)

func initDBConnection() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if err := db.HandlerConnector.Open(dbHost, dbPort, dbUser, dbPassword, dbName); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := initDBConnection(); err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	routes.SetHandlerFunc(routes.InitRouter())

	log.Println("Server listening on http://localhost:5000")

	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal(err)
	}

	defer func(Handler db.Connector) {
		err := Handler.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db.HandlerConnector)

}
