package chat

import "time"

type Message struct {
	ID        	int       	`json:"id"`
	SenderID  	int32       `json:"sender_id"`
	ReceiverID 	int32      	`json:"receiver_id"`
	Content   	string    	`json:"content"`
	Delivered 	bool 		`json:"delivered"`
	CreatedAt  	time.Time 	`json:"created_at"`
}

type IncomingMessage struct {
	ReceiverID 	int32      	`json:"receiver_id"`
	Content   	string    	`json:"content"`
}

type SaveMessageParams struct {
	SenderID  	int32       `json:"sender_id"`
	ReceiverID 	int32      	`json:"receiver_id"`
	Content   	string    	`json:"content"`
}