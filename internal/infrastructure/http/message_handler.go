package http

import (
	"encoding/json"
	"net/http"

	"github.com/HugoBritez/utic.dev-server/internal/application/messages"
	"github.com/HugoBritez/utic.dev-server/internal/domain/repositories"
)

type MessageHandler struct {
	repo          repositories.MessageRepository
	createUseCase *messages.CreateMessageUseCase
}

func NewMessageHandler(repo repositories.MessageRepository, createUseCase *messages.CreateMessageUseCase) *MessageHandler {
	return &MessageHandler{
		repo:          repo,
		createUseCase: createUseCase,
	}
}

type createMessageRequest struct {
	SenderName        string `json:"sender_name"`
	SenderPhoneNumber string `json:"sender_phone_number"`
	MessageText       string `json:"message_text"`
}

func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var req createMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.SenderName == "" || req.SenderPhoneNumber == "" || req.MessageText == "" {
		http.Error(w, `{"error": "sender_name, sender_phone_number and message_text are required"}`, http.StatusBadRequest)
		return
	}

	message, err := h.createUseCase.Execute(r.Context(), req.SenderName, req.SenderPhoneNumber, req.MessageText)
	if err != nil {
		http.Error(w, `{"error": "failed to create message"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := h.repo.GetMessages(r.Context())
	if err != nil {
		http.Error(w, `{"error": "failed to get messages"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}

func (h *MessageHandler) GetNewPhoneNumbers(w http.ResponseWriter, r *http.Request) {
	phoneNumbers, err := h.repo.GetNewPhoneNumbers(r.Context())
	if err != nil {
		http.Error(w, `{"error": "failed to get phone numbers"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(phoneNumbers)
}
