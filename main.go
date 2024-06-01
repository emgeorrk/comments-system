package main

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	graphqlschema "graphql-comments/graphql"
	"graphql-comments/storage"
	"net/http"
)

func main() {
	var flag int = 1
	if flag == 1 {
		storage.DataBase = storage.NewInMemoryStore()
	} else {
		storage.DataBase, _ = storage.NewPostgresDataStore("postgres://postgres:password@localhost:5432/graphql_comments?sslmode=disable")
	}
	
	var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query:    graphqlschema.QueryType,
		Mutation: graphqlschema.MutationType,
		Types:    []graphql.Type{graphqlschema.PostType, graphqlschema.CommentType},
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
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
