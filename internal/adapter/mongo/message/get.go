package message

import (
	"context"
	"fmt"
	"github.com/ziliscite/messaging-app/internal/adapter/mongo"
	domain "github.com/ziliscite/messaging-app/internal/core/domain/message"
	"go.elastic.co/apm"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *Repository) GetAll(ctx context.Context) (*[]domain.Message, error) {
	span, spanCtx := apm.StartSpan(ctx, "get history", "repository")
	defer span.End()

	var messages []domain.Message

	cursor, err := r.col.Find(spanCtx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("%w: message", mongo.ErrNotFound)
	}

	for cursor.Next(spanCtx) {
		var payload domain.Message
		if err = cursor.Decode(&payload); err != nil {
			return &messages, fmt.Errorf("%w: %v", mongo.ErrDecode, err)
		}

		messages = append(messages, payload)
	}

	return &messages, nil
}
