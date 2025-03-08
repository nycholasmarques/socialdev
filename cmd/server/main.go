package main

import (
	"fmt"
	"log"
	"net/http"

	configs "github.com/nycholasmarques/socialdev/config"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/routes"
	"github.com/rs/cors"
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

	// configuração cors
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Permite todas as origens
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"application/json", "Content-Type"},
	})

	handler := c.Handler(r)

	// Iniciar o servidor
	fmt.Println("Server working on http://localhost:8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}