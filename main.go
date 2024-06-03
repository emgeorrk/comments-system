package main

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"graphql-comments/graphql"
	"graphql-comments/storage"
	"graphql-comments/storage/in-memory"
	"graphql-comments/storage/postgres"
	"log"
	"net/http"
	"os"
)

func main() {
	flag := os.Getenv("STORAGE_TYPE")

	switch flag {
	case "in-memory":
		log.Println("Using in-memory storage")
		storage.DataBase = inMemory.NewInMemoryStore()
	case "postgres":
		log.Println("Using PostgreSQL storage")

		host := os.Getenv("POSTGRES_HOST")
		port := os.Getenv("POSTGRES_PORT")
		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")
		dbname := os.Getenv("POSTGRES_DATABASE")

		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

		var err error
		storage.DataBase, err = postgres.NewPostgresDataStore(psqlInfo)
		if err != nil {
			log.Println("Error connecting to PostgreSQL: ", err)
			return
		}

		log.Println("Successfully connected to PostgreSQL!")
	default:
		log.Fatalf("Unknown storage type: \"%s\"\n", flag)
	}

	var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query:    gql.QueryType,
		Mutation: gql.MutationType,
		Types:    []graphql.Type{gql.PostType, gql.CommentType},
	})

	// Создаем GraphQL обработчик
	graphqlHandler := handler.New(&handler.Config{
		Schema:     &schema,
		Playground: true,
	})

	// Устанавливаем обработчик GraphQL
	http.Handle("/graphql", graphqlHandler)

	// Запускаем сервер на порту 8080
	log.Println("Server is running at http://localhost:8084/graphql")
	err := http.ListenAndServe(":8084", nil)
	if err != nil {
		log.Fatal(err)
	}
}
