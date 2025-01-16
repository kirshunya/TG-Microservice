package model

type User struct {
	ID       int64  `json:"ID"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}
