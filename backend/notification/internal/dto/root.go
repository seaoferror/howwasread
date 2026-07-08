package dto

type RegisterNotificationRequest struct {
	OS              string `json:"os"`
	DevicePushToken string `json:"devicePushToken"`
}
