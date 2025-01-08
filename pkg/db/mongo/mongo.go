package mongo

import (
	"context"
	"github.com/ziliscite/messaging-app/pkg/must"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(connection string) *mongo.Client {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	return must.Must(mongo.Connect(context.Background(), options.Client().ApplyURI(connection).SetServerAPIOptions(serverAPI)))
}
