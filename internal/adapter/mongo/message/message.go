package message

import "go.mongodb.org/mongo-driver/mongo"

type Repository struct {
	col *mongo.Collection
}

func New(client *mongo.Client) *Repository {
	return &Repository{
		col: client.Database("message").Collection("message-history"),
	}
}
