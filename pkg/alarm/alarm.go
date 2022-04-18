package alarm

import (
	"log"
	"time"

	"go-alarm/pkg/other"
)

type Alarm struct {
	stop chan struct{}
}

func NewAlarm() *Alarm {
	return &Alarm{}
}

func (a *Alarm) SetUp(t time.Time, f func()) {
	go func() {
		a.stop = make(chan struct{})
		for {
			select {
			case <-time.Tick(time.Second):
				timeNow := time.Now()

				log.Println(timeNow)

				if t.Hour() == timeNow.Hour() && t.Minute() == timeNow.Minute() {
					f()
					return
				}
			case <-a.stop:
				return
			}
		}
	}()
}

func (a Alarm) Reset() {
	if a.stop == nil {
		return
	}

	if !other.ChannelIsClosed(a.stop) {
		close(a.stop)
	}
}
