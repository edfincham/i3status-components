package components

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type DateWidget struct {
	BaseWidget
}

func NewDateWidget(instance int) *DateWidget {
	w := DateWidget{
		*NewBaseWidget(instance, nil),
	}
	return &w
}

func (w *DateWidget) basicLoop() {
	prefix := "\uf017"

	msg := NewMessage()
	msg.Name = "date"
	msg.Colour = YELLOW
	msg.Instance = strconv.Itoa(w.Instance)
	const layout = "Mon 02/01 15:04"
	for {
		msg.FullText = fmt.Sprintf("%s %s", prefix, time.Now().Format(layout))
		w.Output <- *msg
		time.Sleep(w.Refresh)
	}
}

func (w *DateWidget) readLoop() {
	for {
		i := <-w.Input
		if i.Name == "date" {
			executeTemporaryScript()
			go w.basicLoop()
		}
	}
}

func (w *DateWidget) Start() {
	go w.basicLoop()
	go w.readLoop()
}

func executeTemporaryScript() {
	script := `#!/bin/bash
ncal -y         # displays the year calendar
read -n 1 -r -s # wait for a touch key to exit the terminal
exit
`

	tmpFile, err := ioutil.TempFile("", "i3-sc-date-*.sh")
	if err != nil {
		log.Println("Error creating temporary file:", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(script)); err != nil {
		log.Println("Error writing to temporary file:", err)
		return
	}
	tmpFile.Close()

	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		log.Println("Error making script executable:", err)
		return
	}

	cmd := exec.Command("gnome-terminal", "-x", tmpFile.Name())
	cmd.Start()
	cmd.Wait()
}
