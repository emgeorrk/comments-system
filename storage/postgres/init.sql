CREATE DATABASE graphql_comments;

CREATE TABLE IF NOT EXISTS Posts (
    id VARCHAR(128) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS Comments (
    id VARCHAR(128) PRIMARY KEY,
    post_id VARCHAR(128) NOT NULL,
    parent_comment_id VARCHAR(128),
    content VARCHAR(2000) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES Posts(id),
    FOREIGN KEY (parent_comment_id) REFERENCES Comments(id)
);
