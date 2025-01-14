package ydbmodels

import "github.com/upikoth/aireader-go/internal/models"

type Voice struct {
	Name   string
	Source string
}

func NewYDBVoiceModel(voice *models.Voice) *Voice {
	return &Voice{
		Name:   voice.Name,
		Source: string(voice.Source),
	}
}

func (u *Voice) FromYDBModel() *models.Voice {
	return &models.Voice{
		Name:   u.Name,
		Source: models.VoiceSource(u.Source),
	}
}
