package dto

import "github.com/google/uuid"

type SetDevicePushTokenRequest struct {
	DeviceId        uuid.UUID `json:"deviceId"`
	OS              string    `json:"os"`
	DevicePushToken string    `json:"devicePushToken"`
}
