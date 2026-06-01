package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"dev/task-management/pkg/cache"

	"github.com/google/uuid"
)

// NotificationSubscriber lắng nghe Redis Pub/Sub và forward message đến WebSocket client
type NotificationSubscriber struct {
	hub         *HubConnectionWS
	redisClient *cache.RedisCacheService
}

func NewNotificationSubscriber(hub *HubConnectionWS, redisClient *cache.RedisCacheService) *NotificationSubscriber {
	return &NotificationSubscriber{
		hub:         hub,
		redisClient: redisClient,
	}
}

// SubscribeForUser đăng ký lắng nghe channel `notify:{userId}` cho user vừa kết nối WS.
// Chạy dưới dạng goroutine, tự dừng khi ctx bị cancel (client disconnect).
func (s *NotificationSubscriber) SubscribeAssignTaskForUser(ctx context.Context, userId uuid.UUID) {
	channel := fmt.Sprintf("notify:assign:%s", userId.String())
	pubsub := s.redisClient.Subscribe(ctx, channel)
	defer pubsub.Close()

	log.Printf("Subscribed to Redis channel: %s\n", channel)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Unsubscribed from Redis channel: %s\n", channel)
			return
		case msg, ok := <-pubsub.Channel():
			if !ok {
				return
			}
			var event Event
			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				log.Println("Failed to unmarshal notification event:", err)
				continue
			}
			if err := s.hub.SendEventToUser(userId, &event); err != nil {
				log.Println("Failed to send event to WS user:", err)
			}
		}
	}
}

func (s *NotificationSubscriber) SubscribeUpdateStatusForUser(ctx context.Context, userId uuid.UUID) {
	channel := fmt.Sprintf("notify:status:%s", userId.String())
	pubsub := s.redisClient.Subscribe(ctx, channel)
	defer pubsub.Close()

	log.Printf("Subscribed to Redis channel: %s\n", channel)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Unsubscribed from Redis channel: %s\n", channel)
			return
		case msg, ok := <-pubsub.Channel():
			if !ok {
				return
			}
			var event Event
			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				log.Println("Failed to unmarshal notification event:", err)
				continue
			}
			if err := s.hub.SendEventToUser(userId, &event); err != nil {
				log.Println("Failed to send event to WS user:", err)
			}
		}
	}
}
