package models

type VoiceSource string

const VoiceSourceYandex VoiceSource = "yandex"

type Voice struct {
	Name   string
	Source VoiceSource
}
