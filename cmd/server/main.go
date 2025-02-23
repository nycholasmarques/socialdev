package main

import (
	"fmt"
	"log"
	"net/http"

	configs "github.com/nycholasmarques/socialdev/config"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/routes"
)

func main() {
	// Conectar ao banco de dados
	conn, err := configs.NewConnectionDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	// Criar as queries usando o sqlc
	queries := db.New(conn)

	// Configurar as rotas
	r := routes.SetupRoutes(queries)

	// Iniciar o servidor
	fmt.Println("Server working on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}