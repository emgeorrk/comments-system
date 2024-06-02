package main

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/joho/godotenv"
	"graphql-comments/graphql"
	"graphql-comments/storage"
	"graphql-comments/storage/in-memory"
	"graphql-comments/storage/postgres"
	"log"
	"net/http"
)

func main() {
	envFile, _ := godotenv.Read(".env")
	flag := envFile["STORAGE_TYPE"]

	switch flag {
	case "in-memory":
		log.Println("Using in-memory storage")
		storage.DataBase = inMemory.NewInMemoryStore()
	case "postgres":
		log.Println("Using PostgreSQL storage")

		host := envFile["POSTGRES_HOST"]
		port := envFile["POSTGRES_PORT"]
		user := envFile["POSTGRES_USER"]
		password := envFile["POSTGRES_PASSWORD"]
		dbname := envFile["POSTGRES_DATABASE"]

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
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// Устанавливаем обработчик GraphQL
	http.Handle("/graphql", graphqlHandler)

	// Запускаем сервер на порту 8080
	fmt.Println("Server is running at http://localhost:8080/graphql")
	http.ListenAndServe(":8080", nil)
}
