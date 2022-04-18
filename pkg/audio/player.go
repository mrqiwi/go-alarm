package audio

import (
	"log"
	"os"
	"time"

	"go-alarm/pkg/other"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type Player struct {
	initialized bool
	stop        chan struct{}
}

func NewPlayer() *Player {
	return &Player{}
}

func (p *Player) Play(fileName string) {
	go func() {
		p.stop = make(chan struct{})

		f, err := os.Open(fileName)
		if err != nil {
			log.Println(err)
			return
		}

		s, format, err := mp3.Decode(f)
		if err != nil {
			log.Println(err)
			return
		}
		defer s.Close()

		if !p.initialized {
			err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
			if err != nil {
				log.Println(err)
				return
			}
			p.initialized = true
		}

		speaker.Play(beep.Seq(s, beep.Callback(p.Stop)))

		<-p.stop
		speaker.Clear()
	}()
}

func (p Player) Stop() {
	if p.stop == nil {
		return
	}

	if !other.ChannelIsClosed(p.stop) {
		close(p.stop)
	}
}
