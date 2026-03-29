package repository

import (
	"context"
	"log/slog"
)

func (r *Repository) SetServerIP(ctx context.Context, memberId, ip string) error {
	result := r.client.Do(ctx, r.client.B().Set().Key(memberId).Value(ip).Build())
	if result.Error() != nil {
		slog.Error("fail to save member ip", "err", result.Error())
		return result.Error()
	}
	return nil
}

func (r *Repository) RemoveServerIP(ctx context.Context, memberId string) error {
	result := r.client.Do(ctx, r.client.B().Del().Key(memberId).Build())
	if result.Error() != nil {
		slog.Error("fail to remove member ip", "err", result.Error())
		return result.Error()
	}
	return nil
}
