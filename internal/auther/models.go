package auther

import "github.com/google/uuid"

type AccessData struct {
	Token  string `json:"token,omitempty" binding:"required"`
	ApiUrl string `json:"api_url,omitempty" binding:"required"`
	Path   string `json:"path,omitempty" binding:"required"`
	Method string `json:"method,omitempty" binding:"required"`
}

type Api struct {
	ID  uuid.UUID
	Url string
}

type Client struct {
	ID   uuid.UUID
	Name string
}

type Route struct {
	ID     uuid.UUID
	Method string
	Path   string
}
