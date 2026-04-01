package dto

type SetDevicePushTokenRequest struct {
	OS    string `json:"os"`
	Token string `json:"token"`
}
