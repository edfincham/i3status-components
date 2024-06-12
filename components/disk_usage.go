package components

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

type DiskUsageInfo struct {
	DiskPercent float64
}

type DiskUsageWidget struct {
	BaseWidget
}

func NewDiskUsageWidget(instance int) *DiskUsageWidget {
	w := DiskUsageWidget{
		BaseWidget: *NewBaseWidget(instance, nil),
	}
	return &w
}

func (w *DiskUsageWidget) getStatus() (string, string) {
	di, _ := readDiskUsageInfo()

	var prefix, colour string
	prefix = "\uf1c0"

	if di.DiskPercent < 50 {
		colour = GREEN
	} else if di.DiskPercent < 80 {
		colour = YELLOW
	} else {
		colour = RED
	}
	return fmt.Sprintf("%s %0.1f%%", prefix, di.DiskPercent), colour
}

func (w *DiskUsageWidget) basicLoop() {
	msg := NewMessage()
	msg.Name = "disk"
	msg.Border = BACKGROUND
	msg.BorderRight = 10
	msg.Instance = strconv.Itoa(w.Instance)
	for {
		msg.FullText, msg.Colour = w.getStatus()
		w.Output <- *msg
		time.Sleep(w.Refresh)
	}
}

func (w *DiskUsageWidget) readLoop() {
	for {
		i := <-w.Input
		if i.Name == "disk" {
			cmd := exec.Command("gnome-terminal", "-x", "ncdu")

			cmd.Start()
			cmd.Wait()
			go w.basicLoop()
		}
	}
}

func (w *DiskUsageWidget) Start() {
	go w.basicLoop()
	go w.readLoop()
}

func readDiskUsageInfo() (*DiskUsageInfo, error) {
	diskUsageInfo := new(DiskUsageInfo)

	var stat syscall.Statfs_t
	err := syscall.Statfs("/", &stat)
	if err != nil {
		log.Printf("Error: %v", err)
		return diskUsageInfo, err
	}

	totalBlocks := stat.Blocks * uint64(stat.Bsize)
	freeBlocks := stat.Bfree * uint64(stat.Bsize)
	usedBlocks := totalBlocks - freeBlocks
	percentage := (float64(usedBlocks) / float64(totalBlocks)) * 100.0

	diskUsageInfo.DiskPercent = percentage

	return diskUsageInfo, nil
}
