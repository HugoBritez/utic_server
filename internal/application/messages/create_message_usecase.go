package messages

import (
	"context"

	"github.com/HugoBritez/utic.dev-server/internal/domain/entities"
	"github.com/HugoBritez/utic.dev-server/internal/domain/repositories"
)

type CreateMessageUseCase struct {
	repo repositories.MessageRepository
}

func NewCreateMessageUseCase(repo repositories.MessageRepository) *CreateMessageUseCase {
	return &CreateMessageUseCase{
		repo: repo,
	}
}

func (uc *CreateMessageUseCase) Execute(ctx context.Context, senderName, senderPhoneNumber, messageText string) (*entities.Message, error) {
	message := entities.NewMessage(senderName, senderPhoneNumber, messageText)
	me, err := uc.repo.CreateMessage(ctx, *message)

	if err != nil {
		return nil, err
	}

	return me, nil
}
