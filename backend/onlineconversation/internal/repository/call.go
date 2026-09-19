package repository

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (r *Repository) FindParticipantIds(ctx context.Context, conversationId string) ([]string, error) {
	result := r.valkeyClient.Do(ctx, r.valkeyClient.B().Smembers().Key(conversationId).Build())
	if result.Error() != nil {
		slog.Error("fail to get member ip", "err", result.Error())
		return nil, result.Error()
	}
	value, err := result.AsStrSlice()
	if err != nil {
		slog.Error("fail to get ips value string slice", "err", err)
	}
	return value, nil
}

func (r *Repository) AddParticipantId(ctx context.Context, conversationId string, memberId uuid.UUID) error {
	result := r.valkeyClient.Do(ctx, r.valkeyClient.B().Sadd().Key(conversationId).Member(string(memberId[:])).Build())
	if result.Error() != nil {
		slog.Error("fail to add participant member id to conversation",
			"err", result.Error(), "conversationId", conversationId, "memberId", memberId)
		return result.Error()
	}
	return nil
}

func (r *Repository) RemoveParticipantId(ctx context.Context, conversationId string, memberId uuid.UUID) error {

	result := r.valkeyClient.Do(ctx, r.valkeyClient.B().Srem().Key(conversationId).Member(string(memberId[:])).Build())
	if result.Error() != nil {
		slog.Error("fail to remove participant member id to conversation",
			"err", result.Error(), "conversationId", conversationId, "memberId", memberId)
		return result.Error()
	}
	return nil
}

func (r *Repository) SetServerIP(ctx context.Context, memberId, ip string) error {
	result := r.valkeyClient.Do(ctx, r.valkeyClient.B().Set().Key("conversation"+memberId).Value(ip).Build())
	if result.Error() != nil {
		slog.Error("fail to save member ip", "err", result.Error())
		return result.Error()
	}
	return nil
}

func (r *Repository) RemoveServerIP(ctx context.Context, memberId string) error {
	result := r.valkeyClient.Do(ctx, r.valkeyClient.B().Del().Key("conversation"+memberId).Build())
	if result.Error() != nil {
		slog.Error("fail to remove member ip", "err", result.Error())
		return result.Error()
	}
	return nil
}
