package mstts

import (
	"encoding/json"
	"net/http"
	"sync"
)

const (
	voicesListURL = "https://eastus.api.speech.microsoft.com/cognitiveservices/voices/list"
)

var (
	voiceManagerClient = &http.Client{}
	getVoiceHeader     http.Header
	headerOnce         = &sync.Once{}
)

type Voice struct {
	Name            string `json:"Name"`
	DisplayName     string `json:"DisplayName"`
	LocalName       string `json:"LocalName"`
	ShortName       string `json:"ShortName"`
	Gender          string `json:"Gender"`
	Locale          string `json:"Locale"`
	LocaleName      string `json:"LocaleName"`
	SampleRateHertz string `json:"SampleRateHertz"`
	VoiceType       string `json:"VoiceType"`
	Status          string `json:"Status"`
	WordsPerMinute  string `json:"WordsPerMinute"`
}

type VoiceManager struct {
}

func makeVoiceListRequestHeader() http.Header {
	header := make(http.Header)
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 Edg/107.0.1418.26")
	header.Set("X-Ms-Useragent", "SpeechStudio/2021.05.001")
	header.Set("Content-Type", "application/json")
	header.Set("Origin", "https://azure.microsoft.com")
	header.Set("Referer", "https://azure.microsoft.com")
	return header
}

func NewVoiceManager() *VoiceManager {
	headerOnce.Do(func() {
		getVoiceHeader = makeVoiceListRequestHeader()
	})
	return &VoiceManager{}
}

func (m *VoiceManager) ListVoices() ([]Voice, error) {
	req, err := http.NewRequest("GET", voicesListURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header = getVoiceHeader
	resp, err := voiceManagerClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []Voice
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
