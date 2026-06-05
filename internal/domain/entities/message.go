package entities

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID                string    `json:"id"`
	SenderName        string    `json:"sender_name"`
	SenderPhoneNumber string    `json:"sender_phone_number"`
	MessageText       string    `json:"message_text"`
	CreatedAt         time.Time `json:"created_at"`
}

func NewMessage(senderName, senderPhoneNumber, messageText string) *Message {
	return &Message{
		ID:                uuid.New().String(),
		SenderName:        senderName,
		SenderPhoneNumber: senderPhoneNumber,
		MessageText:       messageText,
		CreatedAt:         time.Now(),
	}
}
