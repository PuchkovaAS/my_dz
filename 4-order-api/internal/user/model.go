package user

import "gorm.io/gorm"

type UserId struct {
	gorm.Model
	Phone     string
	SessionId string
	Code      uint
}
