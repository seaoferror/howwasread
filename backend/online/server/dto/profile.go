package dto

type SetNameRequest struct {
	Name string `json:"name"`
}

type GetMyProfileResponse struct {
	Name string `json:"name"`
}
