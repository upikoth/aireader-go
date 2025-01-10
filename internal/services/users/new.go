package users

import (
	"github.com/google/uuid"
	"github.com/upikoth/aireader-go/internal/models"
)

type Option func(user *models.User)

func newUser(
	email string,
	options ...Option,
) *models.User {
	user := &models.User{
		ID:    models.UserID(uuid.New().String()),
		Email: email,
		Role:  models.UserRoleUser,
	}

	for _, option := range options {
		option(user)
	}

	return user
}

func withPasswordHash(ph string) Option {
	return func(user *models.User) {
		user.PasswordHash = ph
	}
}

func withVkID(vkID string) Option {
	return func(user *models.User) {
		user.VkID = vkID
	}
}

func withMailRuID(mailRuID string) Option {
	return func(user *models.User) {
		user.MailRuID = mailRuID
	}
}

func withYandexID(yandexID string) Option {
	return func(user *models.User) {
		user.YandexID = yandexID
	}
}
