package components

import (
	"fmt"
	"os/exec"
	"strconv"
)

type LogoutWidget struct {
	BaseWidget
}

func NewLogoutWidget(instance int) *LogoutWidget {
	w := LogoutWidget{
		*NewBaseWidget(instance, nil),
	}
	return &w
}

func (w *LogoutWidget) message() string {
	return "\uf011"
}

func (w *LogoutWidget) basicLoop() {
	msg := NewMessage()
	msg.Name = "logout"
	msg.Colour = WHITE
	msg.Instance = strconv.Itoa(w.Instance)
	msg.FullText = fmt.Sprint(w.message())
	w.Output <- *msg
}

func (w *LogoutWidget) readLoop() {
	for {
		i := <-w.Input
		if i.Name == "logout" {
			cmd := exec.Command("i3-nagbar", "-t", "warning", "-m", "Log out?", "-b", "Yes", "i3-msg", "exit")

			cmd.Start()
			cmd.Wait()

			go w.basicLoop()
		}
	}
}

func (w *LogoutWidget) Start() {
	go w.basicLoop()
	go w.readLoop()
}
