package components

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type VolumeWidget struct {
	BaseWidget
}

type VolumeInfo struct {
	VolumePercent int
	IsMuted       bool
}

func NewVolumeWidget(instance int, refresh *time.Duration) *VolumeWidget {
	w := VolumeWidget{
		*NewBaseWidget(instance, refresh),
	}
	return &w
}

func (w *VolumeWidget) basicLoop() {
	msg := NewMessage()
	msg.Name = "volume"
	msg.Colour = WHITE
	msg.Instance = strconv.Itoa(w.Instance)
	for {
		msg.FullText, msg.Colour = w.getStatus()
		w.Output <- *msg
		time.Sleep(w.Refresh)
	}
}

func (w *VolumeWidget) readLoop() {
	for {
		i := <-w.Input
		if i.Name == "volume" {
			cmd := exec.Command("gnome-terminal", "-x", "alsamixer")
			cmd.Start()
			cmd.Wait()
			go w.basicLoop()
		}
	}
}

func (w *VolumeWidget) Start() {
	go w.basicLoop()
	go w.readLoop()
}

func (w *VolumeWidget) getStatus() (string, string) {
	vi, _ := readVolumeInfo()

	var prefix, colour string

	if vi.IsMuted {
		colour = RED
		prefix = "\uf026"
		return fmt.Sprintf("%s 0%%", prefix), colour
	} else {
		if vi.VolumePercent < 50 {
			colour = YELLOW
			prefix = "\uf027"
		} else {
			colour = GREEN
			prefix = "\uf028"
		}
		return fmt.Sprintf("%s %d%%", prefix, vi.VolumePercent), colour
	}
}

func readVolumeInfo() (*VolumeInfo, error) {
	volumeInfo := new(VolumeInfo)

	cmdAmixer := exec.Command("amixer", "sget", "Master")
	cmdAwk := exec.Command("awk", "-F[][]", "/dB/ { print $2 }{ print $6 }")

	cmdAwk.Stdin, _ = cmdAmixer.StdoutPipe()

	var cmdAwkOutput bytes.Buffer
	cmdAwk.Stdout = &cmdAwkOutput

	if err := cmdAmixer.Start(); err != nil {
		log.Printf("Error starting cmdAmixer: %v", err)
		return volumeInfo, err
	}

	if err := cmdAwk.Start(); err != nil {
		log.Printf("Error starting cmdAwk: %v", err)
		return volumeInfo, err
	}

	if err := cmdAmixer.Wait(); err != nil {
		log.Printf("Error starting cmdAmixer: %v", err)
		return volumeInfo, err
	}

	if err := cmdAwk.Wait(); err != nil {
		log.Printf("Error waiting on cmdAwk: %v", err)
		return volumeInfo, err
	}

	result := strings.Split(strings.TrimSpace(string(cmdAwkOutput.String())), "\n")

	if len(result) != 2 {
		log.Println("Error extracting volume and mute status")
		return volumeInfo, nil
	}

	percentStr := strings.TrimSuffix(result[0], "%")
	percentInt, err := strconv.Atoi(percentStr)
	if err != nil {
		log.Printf("Error converting percentage string to integer: %v", err)
		return volumeInfo, nil
	}
	volumeInfo.VolumePercent = percentInt

	if result[1] == "on" {
		volumeInfo.IsMuted = false
	} else {
		volumeInfo.IsMuted = true
	}

	return volumeInfo, nil
}
