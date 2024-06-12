package components

import (
	"strconv"
	"time"
)

type SeparatorWidget struct {
	BaseWidget
}

func NewSeparatorWidget(instance int) *SeparatorWidget {
	w := SeparatorWidget{
		*NewBaseWidget(instance, nil),
	}
	return &w
}

func (w *SeparatorWidget) basicLoop() {
	msg := NewMessage()
	msg.FullText = "\ue0b3"
	msg.Name = "separator"
	msg.Colour = WHITE
	msg.Instance = strconv.Itoa(w.Instance)
	msg.Separator = false
	msg.Border = BACKGROUND
	msg.BorderTop = 2
	msg.BorderBottom = 2
	msg.BorderLeft = 6
	msg.BorderRight = 6
	for {
		w.Output <- *msg
		time.Sleep(w.Refresh)
	}
}

func (w *SeparatorWidget) Start() {
	go w.basicLoop()
	go w.readLoop()
}
