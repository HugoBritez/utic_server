package repositories

import (
	"context"

	"github.com/HugoBritez/utic.dev-server/internal/domain/entities"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, message entities.Message) (*entities.Message, error)
	GetMessages(ctx context.Context) ([]entities.Message, error)
	GetNewPhoneNumbers(ctx context.Context) ([]string, error)
}
