package repository

import (
	"context"
	"database/sql"

	"github.com/HugoBritez/utic.dev-server/internal/domain/entities"
	"github.com/HugoBritez/utic.dev-server/internal/domain/repositories"
	"github.com/HugoBritez/utic.dev-server/internal/infrastructure/db"
)

type MessageRepository struct {
	queries *db.Queries
	db      *sql.DB
}

func NewMessageRepository(queries *db.Queries, db *sql.DB) repositories.MessageRepository {
	return &MessageRepository{
		queries: queries,
		db:      db,
	}
}

func messageToEntity(m db.Messages) entities.Message {
	return entities.Message{
		ID:                m.ID,
		SenderName:        m.SenderName,
		SenderPhoneNumber: m.SenderPhoneNumber,
		MessageText:       m.MessageText,
		CreatedAt:         m.CreatedAt,
	}
}

func messageToEntityList(messages []db.Messages) []entities.Message {
	result := make([]entities.Message, len(messages))
	for i, m := range messages {
		result[i] = messageToEntity(m)
	}
	return result
}

func (r *MessageRepository) CreateMessage(ctx context.Context, message entities.Message) (*entities.Message, error) {
	m, err := r.queries.CreateMessage(ctx, db.CreateMessageParams{
		ID:                message.ID,
		SenderName:        message.SenderName,
		SenderPhoneNumber: message.SenderPhoneNumber,
		MessageText:       message.MessageText,
		CreatedAt:         message.CreatedAt,
	})

	if err != nil {
		return nil, err
	}

	e := messageToEntity(m)

	return &e, nil
}

func (r *MessageRepository) GetMessages(ctx context.Context) ([]entities.Message, error) {
	messages, err := r.queries.GetMessages(ctx)
	if err != nil {
		return nil, err
	}

	return messageToEntityList(messages), nil
}

func (r *MessageRepository) GetNewPhoneNumbers(ctx context.Context) ([]string, error) {
	phoneNumbers, err := r.queries.GetNewPhoneNumbers(ctx)
	if err != nil {
		return nil, err
	}

	return phoneNumbers, nil
}
