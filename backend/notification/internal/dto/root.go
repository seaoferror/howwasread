package dto

import "github.com/google/uuid"

type RegisterNotificationRequest struct {
	DeviceId        uuid.UUID `json:"deviceId"` //Deprecated
	OS              string    `json:"os"`
	DevicePushToken string    `json:"devicePushToken"`
}
