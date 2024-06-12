package components

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

type CpuInfo struct {
	CpuPercent float64
}

type CpuWidget struct {
	BaseWidget
}

func NewCpuWidget(instance int) *CpuWidget {
	w := CpuWidget{
		BaseWidget: *NewBaseWidget(instance, nil),
	}
	return &w
}

func (w *CpuWidget) getStatus() (string, string) {
	ci, _ := readCpuInfo()

	var prefix, colour string
	prefix = "\uf080"

	if ci.CpuPercent < 50 {
		colour = GREEN
	} else if ci.CpuPercent < 80 {
		colour = YELLOW
	} else {
		colour = RED
	}
	return fmt.Sprintf("%s %0.1f%%", prefix, ci.CpuPercent), colour
}

func (w *CpuWidget) basicLoop() {
	msg := NewMessage()
	msg.Name = "cpu"
	msg.Instance = strconv.Itoa(w.Instance)
	for {
		msg.FullText, msg.Colour = w.getStatus()
		w.Output <- *msg
		time.Sleep(w.Refresh)
	}
}

func (w *CpuWidget) readLoop() {
	for {
		i := <-w.Input
		if i.Name == "cpu" {
			cmd := exec.Command("gnome-terminal", "-x", "htop")

			cmd.Start()
			cmd.Wait()
			go w.basicLoop()
		}
	}
}

func (w *CpuWidget) Start() {
	go w.basicLoop()
	go w.readLoop()
}

func readCpuInfo() (*CpuInfo, error) {
	cpuInfo := new(CpuInfo)

	percentages, err := cpu.Percent(1*time.Second, false)
	if err != nil {
		log.Printf("Error: %v", err)
		return cpuInfo, nil
	}
	cpuInfo.CpuPercent = percentages[0]

	return cpuInfo, nil
}
