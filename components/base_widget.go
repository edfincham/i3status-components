package components

import (
	"strconv"
	"time"
)

var instanceCount int

type Widget interface {
	Start()
	SetChannels(chan Message, chan Entry)
}

type BaseWidget struct {
	Output   chan Message
	Input    chan Entry
	Refresh  time.Duration
	Instance int
}

func (w *BaseWidget) SetChannels(out chan Message, in chan Entry) {
	w.Output = out
	w.Input = in
}

func NewBaseWidget(instance int, refresh *time.Duration) *BaseWidget {
	if refresh == nil {
		refresh = &REFRESH
	}

	w := BaseWidget{
		Output:   nil,
		Input:    nil,
		Refresh:  *refresh,
		Instance: instance,
	}
	return &w
}

func (w *BaseWidget) basicLoop() {
	msg := NewMessage()
	msg.FullText = "Basic Widget"
	msg.Name = "Basic"
	msg.Colour = WHITE
	msg.Instance = strconv.Itoa(w.Instance)
	for {
		w.Output <- *msg
		time.Sleep(w.Refresh)
	}
}

func (w *BaseWidget) readLoop() {
	for {
		<-w.Input
	}
}

func (w *BaseWidget) Start() {
	go w.basicLoop()
	go w.readLoop()
}
