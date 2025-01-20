package connections

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var PgConn *pgconn.PgConn

var MongodbConn *mongo.Database

func Connect(pg_url string) *pgconn.PgConn {
	conn, err := pgconn.Connect(context.Background(), pg_url)
	if err != nil {
		log.Fatalln("Failed to connect to PostgreSQL:", err)
	}
	return conn
}

func CheckDBConnetions() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Unable to load env file")
		log.Fatal(err.Error())
	}

	pg_url := os.Getenv("PG_CONN_URL")
	mongo_url := os.Getenv("MONGO_CONN_URL")
	// Initialize Postgres connection
	conn := Connect(pg_url)
	PgConn = conn

	// Initialize MongoDB connection
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongo_url))
	if err != nil {
		log.Fatalln("Failed to connect to MongoDB:", err)
	}

	// Ping MongoDB to verify connection
	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalln("Failed to ping MongoDB:", err)
	}

	MongodbConn = mongoClient.Database("wal-logs")
}
