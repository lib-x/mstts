package mstts

import "testing"

func TestVoiceManager_ListVoices(t *testing.T) {
	vm := NewVoiceManager()
	voices, err := vm.ListVoices()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(voices)
}
