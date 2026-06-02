package realtime

import (
	"context"
	"dev/task-management/pkg/response"
	"dev/task-management/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSHandler struct {
	hub        *HubConnectionWS
	subscriber *NotificationSubscriber
}

func NewWSHandler(hub *HubConnectionWS, subscriber *NotificationSubscriber) *WSHandler {
	return &WSHandler{
		hub:        hub,
		subscriber: subscriber,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WSHandler) Connect(c *gin.Context) {
	userId, err := utils.GetOwnerId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ResponseError("Internal Server Error", err.Error()))
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := NewClient(conn)

	h.hub.AddClient(client, userId)

	// Tạo context gắn với vòng đời kết nối WS của user này
	ctx, cancel := context.WithCancel(context.Background())

	// Khi client kết nối, bắt đầu lắng nghe Redis channel `notify:{userId}`
	go h.subscriber.SubscribeAssignTaskForUser(ctx, userId)

	go h.subscriber.SubscribeUpdateStatusForUser(ctx, userId)

	go client.WriteLoop()

	// ReadLoop block cho đến khi client disconnect
	client.ReadLoop()

	// Client đã disconnect → cancel subscriber
	cancel()

	h.hub.RemoveClient(client, userId)
}
