package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var database *sql.DB

func Connect() {
	host := getEnv("DB_HOST", "db")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	name := getEnv("DB_NAME", "carro_ideal")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("erro ao abrir conexão com banco: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("erro ao validar conexão com banco: %v", err)
	}

	database = db
	log.Println("Banco conectado com sucesso!")
}

func GetDB() *sql.DB {
	return database
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
