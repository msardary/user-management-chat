package redisx

import (
	"context"
	"fmt"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type MessageCache struct {
	Client *redis.Client
}

type IncomingMessage struct {
	SenderID  	int32       `json:"sender_id"`
	ReceiverID 	int32      	`json:"receiver_id"`
	Content   	string    	`json:"content"`
}

func NewMessageCache(redisClient *redis.Client) *MessageCache {
	return &MessageCache{
		Client: redisClient,
	}
}

func (r *MessageCache) CacheMessage(ctx context.Context, userID int32, message IncomingMessage) error {
	key := fmt.Sprintf("user:messages:%d", userID)

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	pipe := r.Client.Pipeline()
	pipe.LPush(ctx, key, data)
	pipe.LTrim(ctx, key, 0, 49)
	
	pipe.Expire(ctx, key, 72*60*60) // Set expiration to 3 days

	_, err = pipe.Exec(ctx)
	return err
}

func (r *MessageCache) GetRecentMessages(ctx context.Context, userID int32) ([]IncomingMessage, error) {
	key := fmt.Sprintf("user:messages:%d", userID)

	data, err := r.Client.LRange(ctx, key, 0, 49).Result()
	if err != nil {
		return nil, err
	}

	var messages []IncomingMessage
	for _, item := range data {
		var message IncomingMessage
		err = json.Unmarshal([]byte(item), &message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}