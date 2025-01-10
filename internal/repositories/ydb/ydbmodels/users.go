package ydbmodels

import "github.com/upikoth/aireader-go/internal/models"

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         string
	VkID         string
	MailRuID     string
	YandexID     string
}

func NewYDBUserModel(user *models.User) *User {
	return &User{
		ID:           string(user.ID),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         string(user.Role),
		VkID:         user.VkID,
		MailRuID:     user.MailRuID,
		YandexID:     user.YandexID,
	}
}

func (u *User) FromYDBModel() *models.User {
	return &models.User{
		ID:           models.UserID(u.ID),
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         models.UserRole(u.Role),
		VkID:         u.VkID,
		MailRuID:     u.MailRuID,
		YandexID:     u.YandexID,
	}
}
