package models

type OauthSource string

const OauthSourceVk OauthSource = "vk"
const OauthSourceMail OauthSource = "mail"
const OauthSourceYandex OauthSource = "yandex"

type OauthUserInfo struct {
	ID    string
	Email string
}
