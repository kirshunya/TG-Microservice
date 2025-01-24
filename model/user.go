package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Birthday string `json:"birthday"`
	Age      uint   `json:"age"`
	Link     string `json:"link" gorm:"unique"`
	About    string `json:"about"`
	Active   bool   `json:"active"`
	TgID     uint64 `json:"tg_id" gorm:"unique"`
}
