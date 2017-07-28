package common

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

type Process interface {
	Stop()
	StoppedEvent() chan struct{}
	Run()
	SetAlive(alive bool)
}

func StartProcess(p Process) {
	p.SetAlive(true)
	go p.Run()
}

type BackgroundProcess struct {
	Alive       bool
	ProcessName string
	Done        chan struct{}
	doneLock    sync.RWMutex
	Stopped     chan struct{}
	stoppedLock sync.RWMutex

	Wg sync.WaitGroup
}

func (b *BackgroundProcess) Stop() {
	log.Debugf("Stopping BackgroundProcess: %s", b.ProcessName)

	b.closeDone()
	b.Wg.Wait()
	b.Alive = false
	b.closeStopped()

	log.Debugf("Stopped BackgroundProcess: %s", b.ProcessName)
}

func (b *BackgroundProcess) closeDone() {
	b.doneLock.RLock()
	if !b.IsDone() {
		b.doneLock.RUnlock()
		b.doneLock.Lock()
		close(b.Done)
		b.doneLock.Unlock()
	} else {
		b.doneLock.RUnlock()
	}
}

func (b *BackgroundProcess) closeStopped() {
	b.stoppedLock.RLock()
	if !b.IsStopped() {
		b.stoppedLock.RUnlock()
		b.stoppedLock.Lock()
		close(b.Stopped)
		b.stoppedLock.Unlock()
	} else {
		b.stoppedLock.RUnlock()
	}
}

func (b *BackgroundProcess) IsDone() bool {
	select {
	case <-b.Done:
		return true
	default:
		return false
	}
}

func (b *BackgroundProcess) IsStopped() bool {
	select {
	case <-b.Stopped:
		return true
	default:
		return false
	}
}

func (b *BackgroundProcess) StoppedEvent() chan struct{} {
	return b.Stopped
}

func (b *BackgroundProcess) SetAlive(alive bool) {
	b.Alive = alive
}

func (b *BackgroundProcess) InitProcess(processName string) {
	b.Alive = false
	b.Done = make(chan struct{})
	b.Stopped = make(chan struct{})
	b.ProcessName = processName
}
