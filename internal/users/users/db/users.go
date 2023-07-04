package db

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	iusers "ledger/internal/users/users"
)

type Users struct {
	db     *gorm.DB
	logger *zap.Logger
}
type User struct {
	iusers.User
	Id    string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email string `gorm:"type:varchar(255)"`
}

func NewUsers(dsn string, logger *zap.Logger) (iusers.Users, error) {
	c, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Cannot init Users")
		return nil, err
	}
	u := Users{
		db:     c,
		logger: logger,
	}

	return u, nil
}

func (u Users) GetByEmail(email string) iusers.User {
	user := iusers.User{Email: email}
	//todo check errors
	u.db.First(&user)

	return user
}
