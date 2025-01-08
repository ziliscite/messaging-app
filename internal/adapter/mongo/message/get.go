package message

import (
	"context"
	"fmt"
	"github.com/ziliscite/messaging-app/internal/adapter/mongo"
	domain "github.com/ziliscite/messaging-app/internal/core/domain/message"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *Repository) GetAll(ctx context.Context) (*[]domain.Message, error) {
	var messages []domain.Message

	cursor, err := r.col.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("%w: message", mongo.ErrNotFound)
	}

	for cursor.Next(ctx) {
		var payload domain.Message
		if err = cursor.Decode(&payload); err != nil {
			return &messages, fmt.Errorf("%w: %v", mongo.ErrDecode, err)
		}

		messages = append(messages, payload)
	}

	return &messages, nil
}
