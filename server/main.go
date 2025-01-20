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
	httpPort := fmt.Sprintf(":%s", port)
	router := router.Router()
	router.Run(httpPort)
}
