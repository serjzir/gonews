package postgres

import (
	"GoNews/pkg/storage"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

type Store struct {
	client *pgxpool.Pool
}

func New(connection string) (*Store, error) {
	connect, err := pgxpool.Connect(context.Background(), connection)
	if err != nil {
		return nil, err
	}
	return &Store{client: connect}, nil
}

func (s *Store) Posts() ([]storage.Post, error) {
	var q string
	q = "SELECT posts.id, title, content, posts.author_id, a.name, created_at FROM public.posts LEFT JOIN authors a ON a.id = posts.author_id"
	query, err := s.client.Query(context.Background(), q)
	if err != nil {
		return nil, err
	}
	var posts []storage.Post
	for query.Next() {
		var post storage.Post
		err = query.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.AuthorName, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err = query.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *Store) AddPost(p storage.Post) error {
	q := "INSERT INTO public.posts (author_id, title, content, created_at) VALUES ($1, $2, $3, $4) RETURNING id"
	err := s.client.QueryRow(context.Background(), q, p.AuthorID, p.Title, p.Content, p.CreatedAt).Scan(&p.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) UpdatePost(p storage.Post) error {
	id := p.ID
	q := "UPDATE posts SET title=$1, content=$2 WHERE id=$3"
	_, err := s.client.Exec(context.Background(), q, p.Title, p.Content, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) DeletePost(p storage.Post) error {
	id := p.ID
	q := "DELETE FROM posts WHERE id = $1"
	_, err := s.client.Exec(context.Background(), q, id)
	if err != nil {
		return err
	}
	return nil
}
