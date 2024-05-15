package service

import (
	"context"
	"fmt"
	"github.com/SherClockHolmes/webpush-go"
	"github.com/hashicorp/go-multierror"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
	"pwa/internal/models"
	"pwa/internal/repository"
)

type WebPushService struct {
	repo            *repository.WebPushRepository
	vapidPublicKey  string
	vapidPrivateKey string
	vapidContact    string
	channelRepo     *repository.ChannelRepository
}

func NewWebPushService(repo *repository.WebPushRepository, channelRepo *repository.ChannelRepository) *WebPushService {
	vapidPublicKey := os.Getenv("VAPID_PUBLIC_KEY")
	vapidPrivateKey := os.Getenv("VAPID_PRIVATE_KEY")
	vapidContact := os.Getenv("VAPID_CONTACT")

	if vapidPublicKey == "" || vapidPrivateKey == "" {
		log.Fatal("VAPID keys must be set in the environment")
	}

	return &WebPushService{
		repo:            repo,
		channelRepo:     channelRepo,
		vapidPublicKey:  vapidPublicKey,
		vapidPrivateKey: vapidPrivateKey,
		vapidContact:    vapidContact,
	}
}

func (s *WebPushService) SubscribeUser(ctx context.Context, subscription models.WebPushSubscription) error {
	_, err := s.repo.CreateWebPushSubscription(ctx, subscription)
	return err
}

func (s *WebPushService) NotifyUser(ctx context.Context, userID primitive.ObjectID, message string) error {
	subscriptions, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}
	return s.sendNotifications(ctx, subscriptions, message)
}

func (s *WebPushService) NotifyChannelMembers(ctx context.Context, channelID string, message string) error {
	userIDs, err := s.channelRepo.GetChannelMembers(ctx, channelID)
	if err != nil {
		return fmt.Errorf("failed to get channel members: %w", err)
	}

	var allErrors error
	for _, userIDString := range userIDs {
		userID, err := primitive.ObjectIDFromHex(userIDString)
		if err != nil {
			allErrors = multierror.Append(allErrors, err)
			continue
		}

		subscriptions, err := s.repo.FindByUserID(ctx, userID)
		if err != nil {
			allErrors = multierror.Append(allErrors, fmt.Errorf("failed to find subscriptions for user %s: %w", userIDString, err))
			continue
		}

		if err := s.sendNotifications(ctx, subscriptions, message); err != nil {
			allErrors = multierror.Append(allErrors, fmt.Errorf("failed to send notifications for user %s: %w", userIDString, err))
		}
	}

	return allErrors
}

func (s *WebPushService) sendNotifications(ctx context.Context, subscriptions []models.WebPushSubscription, message string) error {
	var allErrors error
	for _, sub := range subscriptions {
		_, err := webpush.SendNotification([]byte(message), &webpush.Subscription{
			Endpoint: sub.Endpoint,
			Keys: webpush.Keys{
				P256dh: sub.Keys["p256dh"],
				Auth:   sub.Keys["auth"],
			},
		}, &webpush.Options{
			Subscriber:      s.vapidContact,
			VAPIDPublicKey:  s.vapidPublicKey,
			VAPIDPrivateKey: s.vapidPrivateKey,
			TTL:             30,
		})
		if err != nil {
			allErrors = multierror.Append(allErrors, err)
		}
	}
	return allErrors
}
