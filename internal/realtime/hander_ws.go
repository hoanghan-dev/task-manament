package realtime

import (
	"context"
	"dev/task-management/pkg/apperror"
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
		apperror.HandleError(c, err) // 401 Unauthorized (was 500)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := NewClient(conn)

	h.hub.AddClient(client, userId)

	// Create a context tied to this WebSocket connection's lifecycle
	ctx, cancel := context.WithCancel(context.Background())

	// Start listening to Redis channels for this user
	go h.subscriber.SubscribeAssignTaskForUser(ctx, userId)

	go h.subscriber.SubscribeUpdateStatusForUser(ctx, userId)

	go h.subscriber.SubscribeCommentForUser(ctx, userId)

	go client.WriteLoop()

	// ReadLoop blocks until client disconnects
	client.ReadLoop()

	// Client disconnected → cancel subscriber goroutines
	cancel()

	h.hub.RemoveClient(client, userId)
}
