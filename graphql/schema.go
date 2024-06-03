package gql

import (
	"github.com/graphql-go/graphql"
)

// PostType определяет тип постов для GraphQL
var PostType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Post",
	Fields: graphql.Fields{
		"id": &graphql.Field{
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
		"allowComments": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Boolean),
		},
	},
})

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

// QueryType определяет типы запросов для GraphQL
var QueryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"getPosts": &graphql.Field{
			Name:    "getPosts",
			Type:    graphql.NewList(PostType),
			Resolve: getPostsResolver,
		},
		"getPostByID": &graphql.Field{
			Type: PostType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: getPostByIDResolver,
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
			Resolve: getCommentsResolver,
		},
		"getCommentByID": &graphql.Field{
			Type: CommentType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: getCommentByIDResolver,
		},
		"getNumberOfCommentPages": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
			Args: graphql.FieldConfigArgument{
				"postId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: getNumberOfCommentPagesResolver,
		},
		"getReplies": &graphql.Field{
			Type: graphql.NewList(CommentType),
			Args: graphql.FieldConfigArgument{
				"commentID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: getRepliesResolver,
		},
	},
})

// MutationType определяет мутаций для GraphQL
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
			Resolve: addPostResolver,
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
			Resolve: addCommentResolver,
		},
	},
})
