package mongo

import (
	"GoNews/pkg/storage"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName   = "posts"
	collectionName = "posts"
)

type Store struct {
	client *mongo.Client
}

func New(connection string) (*Store, error) {
	mongoOpts := options.Client().ApplyURI(connection)
	client, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return &Store{client: client}, nil
}

func (s *Store) Posts() ([]storage.Post, error) {
	collection := s.client.Database(databaseName).Collection(collectionName)
	cur, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	var posts []storage.Post
	for cur.Next(context.Background()) {
		var post storage.Post
		fmt.Println(post)
		err := cur.Decode(&post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, cur.Err()
}

func (s *Store) AddPost(p storage.Post) error {
	collection := s.client.Database(databaseName).Collection(collectionName)
	_, err := collection.InsertOne(context.Background(), p)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) UpdatePost(p storage.Post) error {
	id := p.ID
	collection := s.client.Database(databaseName).Collection(collectionName)
	_, err := collection.UpdateOne(context.Background(), bson.M{"id": id}, bson.M{
		"$set": bson.M{"title": p.Title, "content": p.Content}})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) DeletePost(p storage.Post) error {
	id := p.ID
	collection := s.client.Database(databaseName).Collection(collectionName)
	_, err := collection.DeleteOne(context.Background(), bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}
