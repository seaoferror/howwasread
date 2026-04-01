package controller

import (
	"context"

	"github.com/google/uuid"
)

func (c *Controller) NotifyMessaging(ctx context.Context, toIds []uuid.UUID, roomId, fromId uuid.UUID, contentType, content string) error {
	err := c.service.NotifyMessaging(ctx, toIds, roomId, fromId, contentType, content)
	if err != nil {
		return err
	}
	return nil
}
