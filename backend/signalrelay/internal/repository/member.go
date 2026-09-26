package repository

import (
	"context"
	"log/slog"
)

func (r *repository) GetServerIP(ctx context.Context, id string) (string, error) {
	result := r.client.Do(ctx, r.client.B().Get().Key("conversation"+id).Build())
	if result.Error() != nil {
		slog.Error("fail to get member ip", "err", result.Error())
		return "", result.Error()
	}
	value, err := result.ToString()
	if err != nil {
		slog.Error("fail to get member ip string", "err", err)
		return "", err
	}
	return value, nil
}

func (r *repository) RemoveServerIP(ctx context.Context, memberId string) error {
	result := r.client.Do(ctx, r.client.B().Del().Key("conversation"+memberId).Build())
	if result.Error() != nil {
		slog.Error("fail to remove member ip", "err", result.Error())
		return result.Error()
	}
	return nil
}
