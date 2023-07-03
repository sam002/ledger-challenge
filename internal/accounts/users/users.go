package users

import "time"

type User struct {
	Id        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email     string `gorm:"type:varchar(255)"`
	IsDealer  bool
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type Users interface {
	GetByEmail(email string) User
}
