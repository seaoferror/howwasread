package dto

type SetNameRequest struct {
	Name string `json:"name"`
}

type GetMyProfileResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
