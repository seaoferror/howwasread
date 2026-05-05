package repository

import (
	"context"
	"log/slog"
)

func (r *Repository) DidNotification(ctx context.Context, messageId, notificationId string) (bool, error) {
	result := r.client.Do(ctx, r.client.B().Sismember().Key("fcm"+messageId).Member(notificationId).Build())
	if result.Error() != nil {
		slog.Error("fail to check file path", "err", result.Error())
		return false, result.Error()
	}
	did, err := result.AsBool()
	if err != nil {
		slog.Error("fail to get result of checking file path", "err", err)
		return false, err
	}
	return did, nil
}

func (r *Repository) MarkNotification(ctx context.Context, messageId, notificationId string) error {
	key := "fcm" + messageId
	result := r.client.Do(ctx, r.client.B().Sadd().Key(key).Member(notificationId).Build())
	if result.Error() != nil {
		slog.Error("fail to save member ip", "err", result.Error())
		return result.Error()
	}
	ttl := r.client.Do(ctx, r.client.B().Expire().Key(key).Seconds(600).Build())
	if ttl.Error() != nil {
		slog.Error("fail to set ttl", "err", ttl.Error())
		return ttl.Error()
	}
	return nil
}
