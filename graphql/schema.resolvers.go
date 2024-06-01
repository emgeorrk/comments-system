package graphql_schema

import (
	"github.com/graphql-go/graphql"
	"graphql-comments/storage"
)

// CommentType определяет тип комментариев для GraphQL
var CommentType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Comment",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"postId": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"parentCommentId": &graphql.Field{
			Type: graphql.ID,
		},
		"content": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"createdAt": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"replies": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(graphql.ID)),
		},
	},
})

// PostType определяет тип постов для GraphQL
var PostType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Post",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Name: "id",
			Type: graphql.NewNonNull(graphql.ID),
		},
		"title": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"content": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"createdAt": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"comments": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(graphql.ID)),
		},
	},
})

// QueryType определяет типы запросов для GraphQL
var QueryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"getPosts": &graphql.Field{
			Name: "getPosts",
			Type: graphql.NewList(PostType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				posts := storage.DataBase.GetPosts()
				return posts, nil
			},
		},
		"getPostByID": &graphql.Field{
			Type: PostType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(string)
				post, err := storage.DataBase.GetPostByID(id)
				if err != nil {
					return nil, err
				}
				return post, nil
			},
		},
		"getCommentByID": &graphql.Field{
			Type: PostType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(string)
				comment, err := storage.DataBase.GetCommentByID(id)
				if err != nil {
					return nil, err
				}
				return comment, nil
			},
		},
		"getComments": &graphql.Field{
			Type: graphql.NewList(CommentType),
			Args: graphql.FieldConfigArgument{
				"postId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"paginationSize": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				postId, _ := params.Args["postId"].(string)
				comments, err := storage.DataBase.GetComments(postId)
				if err != nil {
					return nil, err
				}
				return comments, nil
			},
		},
	},
})

// MutationType определяет мутации для GraphQL
var MutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"addPost": &graphql.Field{
			Type: PostType,
			Args: graphql.FieldConfigArgument{
				"title": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"content": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				title, _ := params.Args["title"].(string)
				content, _ := params.Args["content"].(string)
				newPost := storage.DataBase.AddPost(title, content)
				return newPost, nil
			},
		},
		"addComment": &graphql.Field{
			Type: CommentType,
			Args: graphql.FieldConfigArgument{
				"postId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"parentCommentId": &graphql.ArgumentConfig{
					Type: graphql.ID,
				},
				"content": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				postID, _ := params.Args["postId"].(string)
				parentCommentID, _ := params.Args["parentCommentId"].(string)
				content, _ := params.Args["content"].(string)
				newComment, err := storage.DataBase.AddComment(postID, parentCommentID, content)
				if err != nil {
					return nil, err
				}
				return newComment, nil
			},
		},
	},
})
