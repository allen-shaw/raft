package raft

import (
	"fmt"
	"github.com/AllenShaw19/raft/log"
	"sync"
	"time"
)

type RepeatedTimerTask struct {
	mutex           sync.Mutex
	timer           *time.Timer
	nextDueTime     time.Duration
	timeoutMs       int
	stopped         bool
	running         bool
	destroyed       bool
	invoking        bool
	run             func()
	onDestroy       func()
	adjustTimeoutMs func(timeoutMs int) int
	exitC           chan struct{}
}

func NewRepeatedTimerTask() *RepeatedTimerTask {
	t := &RepeatedTimerTask{
		timeoutMs: 0,
		stopped:   true,
		running:   false,
		destroyed: true,
		invoking:  false,
		exitC:     make(chan struct{}),
	}
	return t
}

func (t *RepeatedTimerTask) SetRun(f func()) {
	t.run = f
}
func (t *RepeatedTimerTask) SetOnDestroy(f func()) {
	t.onDestroy = f
}
func (t *RepeatedTimerTask) SetAdjustTimeoutMs(f func(timeoutMs int) int) {
	t.adjustTimeoutMs = f
}

func (t *RepeatedTimerTask) Init(timeoutMs int) {
	t.timeoutMs = timeoutMs
	t.destroyed = false
	t.stopped = true
	t.running = false
}

func (t *RepeatedTimerTask) Start() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.destroyed || !t.stopped {
		return
	}
	t.stopped = false
	if t.running {
		return
	}
	t.running = true
	t.timer = time.NewTimer(time.Duration(t.timeoutMs) * time.Millisecond)
	go t.start()
}

func (t *RepeatedTimerTask) start() {
	for {
		select {
		case <-t.timer.C:
			t.onTimeout()
		case <-t.exitC:
			log.Info("task exit")
			return
		}
	}
}

func (t *RepeatedTimerTask) RunOnceNow() {
	t.mutex.Lock()
	if t.timer.Stop() {
		t.onTimeout()
		return
	}
	t.mutex.Unlock()
	return
}

func (t *RepeatedTimerTask) Stop() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.stopped {
		return
	}
	t.stopped = true
	if !t.running {
		log.Fatal("check fail, running %v", t.running)
	}
	ok := t.timer.Stop()
	if ok {
		t.running = false
	}
}

func (t *RepeatedTimerTask) Reset(timeoutMs int) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if timeoutMs != 0 {
		t.timeoutMs = timeoutMs
	}
	if t.stopped {
		return
	}
	if !t.running {
		log.Fatal("check fail, running %v", t.running)
	}
	t.schedule()
}
func (t *RepeatedTimerTask) Destroy() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.destroyed {
		return
	}
	t.destroyed = true
	close(t.exitC)
	if !t.running {
		if !t.stopped {
			log.Fatal("check fail, stopped %v", t.stopped)
		}
		t.onDestroy()
		return
	}
	if t.stopped {
		return
	}
	t.stopped = true
	if t.timer.Stop() {
		t.running = false
		t.onDestroy()
	}
}

func (t *RepeatedTimerTask) Describe() string {
	t.mutex.Lock()
	stopped := t.stopped
	running := t.running
	destroyed := t.destroyed
	invoking := t.invoking
	duetime := t.nextDueTime
	timeoutMs := t.timeoutMs
	t.mutex.Unlock()

	description := fmt.Sprintf("timeout(\"%v\"ms)", timeoutMs)
	if destroyed {
		description += " DESTROYED"
	}
	if stopped {
		description += " STOPPED"
	}
	if running {
		if invoking {
			description += " INVOKING"
		} else {
			description += fmt.Sprintf(" SCHEDULING(in \"%v\"ms", duetime)
		}
	}
	return description
}

func (t *RepeatedTimerTask) onTimeout() {
	t.mutex.Lock()
	t.invoking = true
	t.mutex.Unlock()

	t.run()

	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.invoking = false
	if !t.running {
		log.Fatal("check fail, running %v", t.running)
	}
	if t.stopped {
		t.running = false
		if t.destroyed {
			t.onDestroy()
		}
	}

	t.schedule()
}

func (t *RepeatedTimerTask) schedule() {
	t.nextDueTime = time.Duration(t.adjustTimeoutMs(t.timeoutMs)) * time.Millisecond
	if !t.timer.Reset(t.nextDueTime) {
		log.Error("Fail to reset timer")
		t.onTimeout()
	}
}
