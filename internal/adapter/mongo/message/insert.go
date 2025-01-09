package message

import (
	"context"
	"github.com/ziliscite/messaging-app/internal/core/domain/message"
	"go.elastic.co/apm"
)

func (r *Repository) Insert(ctx context.Context, message *message.Message) error {
	span, spanCtx := apm.StartSpan(ctx, "insert message", "repository")
	defer span.End()

	_, err := r.col.InsertOne(spanCtx, message)
	return err
}
