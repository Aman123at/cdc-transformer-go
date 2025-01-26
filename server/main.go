package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Aman123at/cdc-go/connections"
	router "github.com/Aman123at/cdc-go/routes"
	stream "github.com/Aman123at/cdc-go/wal-stream-pg"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Welcome to cdc")
	go func() {
		log.Println("Starting WAL streaming")
		stream.StartStream()
	}()
	go func() {
		startHttpServer()
	}()
	select {}
}

func init() {
	connections.CheckDBConnetions()

	// initiate pg connection pool
	pool, err := connections.NewConnectionPool(3)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Initiated connection pool")
	connections.PgPool = pool
}

func startHttpServer() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Unable to load env file")
		log.Fatal(err.Error())
	}

	port := os.Getenv("HTTP_PORT")
	go_env := os.Getenv("GO_ENV")
	cert_file_path := os.Getenv("CERT_FILE_PATH")
	key_file_path := os.Getenv("KEY_FILE_PATH")
	httpPort := fmt.Sprintf(":%s", port)
	router := router.Router()
	if go_env == "prod" {
		err := router.RunTLS(httpPort, cert_file_path, key_file_path)
		if err != nil {
			log.Fatal("Failed to start server on prod: ", err)
		}
	} else {
		err := router.Run(httpPort)
		if err != nil {
			log.Fatal("Failed to start server: ", err)
		}
	}
}
