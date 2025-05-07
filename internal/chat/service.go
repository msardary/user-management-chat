package chat

import (
	"context"
	"user-management/internal/db/generated"
	"user-management/pkg/redisx"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID int
	Conn   *websocket.Conn
}

type Service struct {
	db 		*db.Queries
	clients map[int]*Client
	cache 	*redisx.MessageCache
}

func NewService(db *db.Queries, cache *redisx.MessageCache) *Service {
	return &Service{
		db: 		db,
		clients: 	make(map[int]*Client),
		cache: 		cache,
	}
}

func (s *Service) AddClient(userID int, conn *websocket.Conn) {
	s.clients[userID] = &Client{
		UserID: userID,
		Conn:   conn,
	}
}

func (s *Service) RemoveClient(userID int) {
	if client, ok := s.clients[userID]; ok {
		client.Conn.Close()
		delete(s.clients, userID)
	}
}

func (s *Service) SaveMessage(ctx context.Context, payload struct {
	SenderID	int32 	`json:"sender_id"`
	ReceiverID  int32 	`json:"receiver_id"`
	Content 	string 	`json:"content"`
}) (*db.Message, error) {

	params := db.InsertMessageParams{
		SenderID:    	payload.SenderID,
		ReceiverID: 	payload.ReceiverID,
		Content: 		payload.Content,
	}
	
	message, err := s.db.InsertMessage(ctx, params)
	if err != nil {
		return nil, err
	}

	chatMsg := redisx.IncomingMessage{
		SenderID:   payload.SenderID,
		ReceiverID: payload.ReceiverID,
		Content:    payload.Content,
	}
	// Cache the message
	s.cache.CacheMessage(ctx, payload.SenderID, chatMsg)
	s.cache.CacheMessage(ctx, payload.ReceiverID, chatMsg)

	return &message, err
}

func (s *Service) MarkAsDelivered(ctx context.Context, messageID int32) error {
	
	return s.db.MarkAsDelivered(ctx, messageID)

}

func (s *Service) GetUndeliveredMessages(ctx context.Context, userID int32) ([]db.Message, error) {

	return s.db.GetUndeliveredMessages(context.Background(), userID)

}

func (s *Service) GetRecentMessages(ctx context.Context, userID int32) ([]redisx.IncomingMessage, error) {
	return s.cache.GetRecentMessages(ctx, userID)
}