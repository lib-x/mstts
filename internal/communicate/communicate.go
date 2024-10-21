package communicate

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	communicateClient = &http.Client{}
)

const (
	endpointURL          = "https://dev.microsofttranslator.com/apps/endpoint?api-version=1.0"
	userAgent            = "okhttp/4.5.0"
	clientVersion        = "4.0.530a 5fe1dc6c"
	userId               = "0f04d16a175c411e"
	homeGeographicRegion = "zh-Hans-CN"
	clientTraceId        = "aab069b9-70a7-4844-a734-96cd78d94be9"
	voiceDecodeKey       = "oik6PdDdMnOXemTbwvMn9de/h9lFnfBaCWbGMMZqqoSaQaqUOqjVGm5NqsmjcBI1x+sS9ugjB55HEJWRiFXYFw=="
	defaultVoiceName     = "zh-CN-XiaoxiaoMultilingualNeural"
	defaultRate          = "0"
	defaultPitch         = "0"
	defaultOutputFormat  = "audio-24khz-48kbitrate-mono-mp3"
	defaultStyle         = "general"
)

type endpoint struct {
	Region string `json:"r"`
	Token  string `json:"t"`
}

type Communicate struct {
	Voice           string
	VoiceLangRegion string
	Pitch           string
	Rate            string
	Volume          string

	outputFormat          string
	style                 string
	IgnoreSSLVerification bool
}

func New(options ...Option) *Communicate {
	c := &Communicate{}
	for _, option := range options {
		option(c)
	}
	if c.Voice == "" {
		c.Voice = defaultVoiceName
	}
	if c.Rate == "" {
		c.Rate = defaultRate
	}
	if c.Pitch == "" {
		c.Pitch = defaultPitch
	}
	if c.outputFormat == "" {
		c.outputFormat = defaultOutputFormat
	}
	if c.style == "" {
		c.style = defaultStyle
	}
	return c
}

func sign(urlStr string) string {
	u := strings.Split(urlStr, "://")[1]
	encodedUrl := url.QueryEscape(u)
	uuidStr := strings.ReplaceAll(uuid.New().String(), "-", "")
	formattedDate := strings.ToLower(time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05")) + "gmt"
	bytesToSign := fmt.Sprintf("MSTranslatorAndroidApp%s%s%s", encodedUrl, formattedDate, uuidStr)
	bytesToSign = strings.ToLower(bytesToSign)
	decode, _ := base64.StdEncoding.DecodeString(voiceDecodeKey)
	hash := hmac.New(sha256.New, decode)
	hash.Write([]byte(bytesToSign))
	secretKey := hash.Sum(nil)
	signBase64 := base64.StdEncoding.EncodeToString(secretKey)
	return fmt.Sprintf("MSTranslatorAndroidApp::%s::%s::%s", signBase64, formattedDate, uuidStr)
}

func (c *Communicate) getEndpoint() (*endpoint, error) {
	signature := sign(endpointURL)
	req, err := http.NewRequest("POST", endpointURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header = makeGetEndpointRequestHeader(signature)
	resp, err := communicateClient.Do(req)
	if err != nil {
		log.Println("failed to do request: ", err)
		return nil, err
	}
	defer resp.Body.Close()
	var result endpoint
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Communicate) GenerateVoiceStreamTo(text string, target io.Writer) error {
	endpoint, err := c.getEndpoint()
	if err != nil {
		return err
	}
	u := fmt.Sprintf("https://%s.tts.speech.microsoft.com/cognitiveservices/v1", endpoint.Region)
	headers := map[string]string{
		"Authorization":            endpoint.Token,
		"Content-Type":             "application/ssml+xml",
		"X-Microsoft-OutputFormat": c.outputFormat,
	}
	ssml := c.getSsml(text)
	req, err := http.NewRequest("POST", u, bytes.NewBufferString(ssml))
	if err != nil {
		return err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := communicateClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(target, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// GenerateVoice generate the voice
func (c *Communicate) GenerateVoice(text string) ([]byte, error) {
	var buf bytes.Buffer
	err := c.GenerateVoiceStreamTo(text, &buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GetSsml 生成 SSML 格式的文本
func (c *Communicate) getSsml(text string) string {
	text = html.EscapeString(text)
	return fmt.Sprintf(`
   <speak xmlns="http://www.w3.org/2001/10/synthesis" xmlns:mstts="http://www.w3.org/2001/mstts" version="1.0" xml:lang="zh-CN">
     <voice name="%s">
       <mstts:express-as style="%s" styledegree="1.0" role="default">
         <prosody rate="%s%%" pitch="%s%%" volume="medium">
			%s
		</prosody>
       </mstts:express-as>
     </voice>
   </speak>
 `, c.Voice, c.style, c.Rate, c.Pitch, text)
}
