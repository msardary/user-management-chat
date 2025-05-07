package chat

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var logger = logrus.StandardLogger()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ChatHandler(service *Service) gin.HandlerFunc {

	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"path":   c.Request.URL.Path,
				"method": c.Request.Method,
				"ip":     c.ClientIP(),
			}).Error("Failed to upgrade connection: ", err)
			log.Println("Failed to upgrade connection: ", err)
			return
		}
		defer conn.Close()

		userIDStr := c.Query("user_id")

		logger.WithFields(logrus.Fields{
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
			"ip":     c.ClientIP(),
		}).Info("userIDStr: ", userIDStr)

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"path":   c.Request.URL.Path,
				"method": c.Request.Method,
				"ip":     c.ClientIP(),
			}).Error("error query: ", err)
			log.Println("error query: ", err)
			return
		}

		service.AddClient(userID, conn)
		defer service.RemoveClient(userID)
		log.Println("Client connected:", userID)

		// undeliveredMsgs, _ := service.GetUndeliveredMessages(context.Background(), int32(userID))
		// for _, message := range undeliveredMsgs {
		// 	err = conn.WriteJSON(message)
		// 	if err != nil {
		// 		logger.WithFields(logrus.Fields{
		// 			"path":   c.Request.URL.Path,
		// 			"method": c.Request.Method,
		// 			"ip":     c.ClientIP(),
		// 		}).Error("Error sending undelivered message: ", err)
		// 		log.Println("Error sending undelivered message:", err)
		// 		break
		// 	}
		// 	service.MarkAsDelivered(context.Background(), message.ID)
		// }

		recentMsgs, err := service.GetRecentMessages(context.Background(), int32(userID))
		if err != nil {
			logger.WithFields(logrus.Fields{
				"path":   c.Request.URL.Path,
				"method": c.Request.Method,
				"ip":     c.ClientIP(),
			}).Error("Error getting recent messages: ", err)
			return
		} else {
			for _, message := range recentMsgs {
				err = conn.WriteJSON(message)
				if err != nil {
					logger.WithFields(logrus.Fields{
						"path":   c.Request.URL.Path,
						"method": c.Request.Method,
						"ip":     c.ClientIP(),
					}).Error("Error sending recent message: ", err)
					break
				}
			}
		}

		for {

			var msg IncomingMessage
			err := conn.ReadJSON(&msg)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"path":   c.Request.URL.Path,
					"method": c.Request.Method,
					"ip":     c.ClientIP(),
				}).Error("Error reading JSON: ", err)
				log.Println("Error reading JSON:", err)
				break
			}

			saveMessageParams := IncomingMessage{
				SenderID:   msg.SenderID,
				ReceiverID: msg.ReceiverID,
				Content:    msg.Content,
			}

			savedMsg, err := service.SaveMessage(context.Background(), saveMessageParams)

			if err != nil {
				logger.WithFields(logrus.Fields{
					"path":   c.Request.URL.Path,
					"method": c.Request.Method,
					"ip":     c.ClientIP(),
				}).Error("Error saving message: ", err)
				log.Println("Error saving message:", err)
				break
			}

			if receiverClient, ok := service.clients[int(msg.ReceiverID)]; ok {
				err = receiverClient.Conn.WriteJSON(msg)
				service.MarkAsDelivered(context.Background(), savedMsg.ID)
				if err != nil {
					logger.WithFields(logrus.Fields{
						"path":   c.Request.URL.Path,
						"method": c.Request.Method,
						"ip":     c.ClientIP(),
					}).Error("Error sending message to receiver: ", err)
					log.Println("Error sending message to receiver:", err)
					break
				}
			}

			conn.WriteJSON(msg)
		}
	}

}