package ttsTask

import (
	"github.com/lib-x/mstts/internal/communicate"
	"io"
	"log"
	"sync"
)

type SingleTask struct {
	C *communicate.Communicate
	// Text to be synthesized
	Text string
	// Output
	Output io.Writer
}

// Start  start a tts single task
func (t *SingleTask) Start(wg *sync.WaitGroup) error {
	defer wg.Done()
	if err := t.C.GenerateVoiceStreamTo(t.Text, t.Output); err != nil {
		return err
	}
	if closer, ok := t.Output.(io.Closer); ok {
		log.Print("ttsTask.Start: close output writer\r\n")
		closer.Close()
	}
	return nil
}
