package ttsTask

import (
	"encoding/json"
	"github.com/lib-x/mstts/internal/communicate"
	"io"
	"log"
	"sync"
)

type PackEntry struct {
	// Text to be synthesized
	Text string
	// Entry name to be packed into a file
	EntryName string
	// EntryOpts defines the options for communicating with the TTS engine.if note set, use the PackTask's CommunicateOpt.
	EntryOpts []communicate.Option
}

type PackTask struct {
	// PackEntryCreator defines the function to create a writer for each entry
	PackEntryCreator func(string) (io.Writer, error)
	// PackOpts defines the options for communicating with the TTS engine
	PackOpts []communicate.Option
	// PackEntries defines the list of entries to be packed into a file
	PackEntries []*PackEntry
	// Output
	Output io.Writer
	// MetaData is the data which will be serialized into a json file,name use the key and value as the key-value pair.
	MetaData []map[string]any
}

func (p *PackTask) Start(wg *sync.WaitGroup) error {
	defer wg.Done()
	for _, entry := range p.PackEntries {
		// for zip file, the entry should be written after creation.
		if err := p.processPackEntry(entry); err != nil {
			continue
		}
	}
	// after all entries are written, write the meta data into a json file. this process is optional.
	// so error is ignored.
	p.writeMetaDataForPack()
	return nil
}

func (p *PackTask) writeMetaDataForPack() {
	if len(p.MetaData) > 0 {
		for _, metaData := range p.MetaData {
			for entryName, entryPayload := range metaData {
				metaEntry, err := p.PackEntryCreator(entryName)
				if err != nil {
					log.Printf("create meta entry writer error:%v \r\n", err)
					continue
				}
				if err = json.NewEncoder(metaEntry).Encode(entryPayload); err != nil {
					log.Printf("write data to meta entry writer error:%v \r\n", err)
					continue
				}
			}
		}
	}
}

func (p *PackTask) processPackEntry(entry *PackEntry) error {
	opt := p.PackOpts
	if entry.EntryOpts != nil {
		opt = entry.EntryOpts
	}
	c := communicate.New(opt...)
	entryWriter, err := p.PackEntryCreator(entry.EntryName)
	err = c.GenerateVoiceStreamTo(entry.Text, entryWriter)
	if err != nil {
		log.Printf("write data to entry writer error:%v \r\n", err)
		return err
	}
	return nil
}
