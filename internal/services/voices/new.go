package voices

import (
	"github.com/upikoth/aireader-go/internal/models"
)

type Option func(voice *models.Voice)

func newVoice(
	name string,
	source models.VoiceSource,
	options ...Option,
) *models.Voice {
	voice := &models.Voice{
		Name:   name,
		Source: source,
	}

	for _, option := range options {
		option(voice)
	}

	return voice
}
