package components

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

type MemoryInfo struct {
	MemoryPercent float64
}

type MemoryWidget struct {
	BaseWidget
}

func NewMemoryWidget(instance int) *MemoryWidget {
	w := MemoryWidget{
		BaseWidget: *NewBaseWidget(instance, nil),
	}
	return &w
}

func (w *MemoryWidget) getStatus() (string, string) {
	mi, _ := readMemoryInfo()

	var prefix, colour string
	prefix = "\uf0a0"

	if mi.MemoryPercent < 50 {
		colour = GREEN
	} else if mi.MemoryPercent < 80 {
		colour = YELLOW
	} else {
		colour = RED
	}
	return fmt.Sprintf("%s %0.1f%%", prefix, mi.MemoryPercent), colour
}

func (w *MemoryWidget) basicLoop() {
	msg := NewMessage()
	msg.Name = "cpu"
	msg.Border = BACKGROUND
	msg.BorderRight = 10
	msg.Instance = strconv.Itoa(w.Instance)
	for {
		msg.FullText, msg.Colour = w.getStatus()
		w.Output <- *msg
		time.Sleep(w.Refresh)
	}
}

func (w *MemoryWidget) readLoop() {
	for {
		i := <-w.Input
		if i.Name == "memory" {
			cmd := exec.Command("gnome-terminal", "-x", "htop")

			cmd.Start()
			cmd.Wait()
			go w.basicLoop()
		}
	}
}

func (w *MemoryWidget) Start() {
	go w.basicLoop()
	go w.readLoop()
}

func readMemoryInfo() (*MemoryInfo, error) {
	memoryInfo := new(MemoryInfo)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	totalMemory := float64(memStats.Sys)
	usedMemory := float64(memStats.HeapInuse + memStats.StackInuse)
	memoryInfo.MemoryPercent = (usedMemory / totalMemory) * 100.0

	return memoryInfo, nil
}
