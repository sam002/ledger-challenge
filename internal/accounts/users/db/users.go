package db

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	iusers "ledger/internal/accounts/users"
)

type users struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUsers(dsn string, logger *zap.Logger) (iusers.Users, error) {
	c, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Cannot init Users")
		return nil, err
	}
	u := users{
		db:     c,
		logger: logger,
	}

	return u, nil
}

func (u users) GetByEmail(email string) iusers.User {
	user := iusers.User{Email: email}
	//todo check errors
	u.db.First(&user)

	return user
}
