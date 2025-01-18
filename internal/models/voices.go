package models

type VoiceSource string

const VoiceSourceYandex VoiceSource = "yandex"

type Voice struct {
	Name   string
	Source VoiceSource
}

type VoiceList struct {
	Voices []*Voice
	Total  int
}
