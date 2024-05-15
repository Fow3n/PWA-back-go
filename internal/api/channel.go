package api

type CreateChannelRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
