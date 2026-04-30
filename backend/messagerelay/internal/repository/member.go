package repository

import (
	"context"
	"log/slog"
)

func (r *Repository) GetServerIPs(ctx context.Context, id string) ([]string, error) {
	result := r.client.Do(ctx, r.client.B().Smembers().Key(id).Build())
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

func (r *Repository) RemoveServerIP(ctx context.Context, memberId, ip string) error {
	result := r.client.Do(ctx, r.client.B().Srem().Key(memberId).Member(ip).Build())
	if result.Error() != nil {
		slog.Error("fail to remove member ip", "err", result.Error())
		return result.Error()
	}
	return nil
}
