package mstts

import "sync"

type Tasker interface {
	Start(wg *sync.WaitGroup) error
}
