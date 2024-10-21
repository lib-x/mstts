package communicate

type Option func(option *Communicate)

func WithVoice(voice string) Option {
	return func(option *Communicate) {
		option.Voice = voice
	}
}

func WithVoiceLangRegion(voiceLangRegion string) Option {
	return func(option *Communicate) {
		option.VoiceLangRegion = voiceLangRegion
	}

}

// WithPitch set pitch of the tts output.such as +50Hz,-50Hz
func WithPitch(pitch string) Option {
	return func(option *Communicate) {
		option.Pitch = pitch
	}
}

// WithRate set rate of the tts output.rate=-50% means rate down 50%,rate=+50% means rate up 50%
func WithRate(rate string) Option {
	return func(option *Communicate) {
		option.Rate = rate
	}
}

// WithVolume set volume of the tts output.volume=-50% means volume down 50%,volume=+50% means volume up 50%
func WithVolume(volume string) Option {
	return func(option *Communicate) {
		option.Volume = volume
	}
}
