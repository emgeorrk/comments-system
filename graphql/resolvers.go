package gql

import (
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	"graphql-comments/storage"
)

func addPostResolver(params graphql.ResolveParams) (interface{}, error) {
	title, _ := params.Args["title"].(string)
	content, _ := params.Args["content"].(string)
	allowComments, ok := params.Args["allowComments"].(bool)

	switch {
	case title == "":
		return nil, errors.New("title is empty")
	case content == "":
		return nil, errors.New("content is empty")
	case len(title) > storage.MaxPostTitleLength:
		return nil, errors.New(fmt.Sprintf("title is too long (maximum %d chars)", storage.MaxPostTitleLength))
	case len(content) > storage.MaxPostContentLength:
		return nil, errors.New(fmt.Sprintf("content is too long (maximum %d chars)", storage.MaxPostContentLength))
	}

	if !ok {
		allowComments = true
	}
	newPost, err := storage.DataBase.AddPost(title, content, allowComments)
	if err != nil {
		return nil, err
	}
	return newPost, nil
}

func addCommentResolver(params graphql.ResolveParams) (interface{}, error) {
	postID, _ := params.Args["postId"].(string)
	parentCommentID, _ := params.Args["parentCommentId"].(string)
	content, _ := params.Args["content"].(string)

	switch {
	case postID == "":
		return nil, errors.New("postID is empty")
	case content == "":
		return nil, errors.New("content is empty")
	case len(content) > storage.MaxCommentLength:
		return nil, errors.New(fmt.Sprintf("content is too long (maximum %d chars)", storage.MaxCommentLength))
	}

	newComment, err := storage.DataBase.AddComment(postID, parentCommentID, content)
	if err != nil {
		return nil, err
	}
	return newComment, nil
}

func getPostsResolver(params graphql.ResolveParams) (interface{}, error) {
	posts, err := storage.DataBase.GetPosts()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func getPostByIDResolver(params graphql.ResolveParams) (interface{}, error) {
	id, _ := params.Args["id"].(string)
	post, err := storage.DataBase.GetPostByID(id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func getCommentsResolver(params graphql.ResolveParams) (interface{}, error) {
	postId, _ := params.Args["postId"].(string)
	comments, err := storage.DataBase.GetComments(postId)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func getCommentByIDResolver(params graphql.ResolveParams) (interface{}, error) {
	id, _ := params.Args["id"].(string)
	comment, err := storage.DataBase.GetCommentByID(id)
	if err != nil {
		return nil, err
	}
	return comment, nil
}
