package message

import (
	"context"
	"github.com/ziliscite/messaging-app/internal/core/domain/message"
)

func (r *Repository) Insert(ctx context.Context, message *message.Message) error {
	_, err := r.col.InsertOne(ctx, message)
	return err
}
