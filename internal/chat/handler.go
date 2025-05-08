package chat

import (
	"context"
	"log"
	"net/http"
	"user-management/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var logger = logrus.StandardLogger()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ChatHandler handles WebSocket connections for real-time chat.
//
// @Summary     WebSocket chat connection
// @Description Connects to WebSocket server for real-time messaging. Token must be passed in query param. This is a WebSocket endpoint.
// @Tags        chat
// @Produce     json
// @Param       token query string true "Access token"
// @Success     101 {string} string "WebSocket upgrade successful"
// @Failure     401 {object} map[string]string "Unauthorized"
// @Failure     500 {object} map[string]string "Internal Server Error"
// @Router      /ws/chat [get]
func ChatHandler(service *Service) gin.HandlerFunc {

	return func(c *gin.Context) {

		tokenString := c.DefaultQuery("token", "")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required in URL!"})
            return
        }

		userID, isAdmin, err := utils.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not invalid!"})
			return
		}

		c.Set("userID", userID)
        c.Set("isAdmin", isAdmin)

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

		service.AddClient(int(userID), conn)
		defer service.RemoveClient(int(userID))

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

			saveMessageParams := SaveMessageParams{
				SenderID:   int32(userID),
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