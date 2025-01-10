package models

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type UserAction string

const (
	UserActionGetAnyUserInfo UserAction = "get-any-user-info"
)

var userRoleActionMap = map[UserRole]map[UserAction]bool{
	UserRoleUser: {
		UserActionGetAnyUserInfo: false,
	},
	UserRoleAdmin: {
		UserActionGetAnyUserInfo: true,
	},
}

func (u UserRole) CheckAccessToAction(action UserAction) bool {
	return userRoleActionMap[u][action]
}

type UserID string

type User struct {
	ID           UserID
	Email        string
	PasswordHash string
	Role         UserRole
	VkID         string
	MailRuID     string
	YandexID     string
}

type UserList struct {
	Users []*User
	Total int
}
