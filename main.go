package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/edfincham/i3status-components/components"
)

func main() {
	debug := os.Getenv("I3STATUS_DEBUG")

	if debug == "true" {
		logFile, err := os.OpenFile("i3status.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("Error opening log file: %v", err)
		}
		defer logFile.Close()

		log.SetOutput(logFile)
	}

	fmt.Println(`{"version":1,"click_events": true}`)
	fmt.Println("[")

	b := components.NewBar()

	powerRefresh := 5000 * time.Millisecond
	volumeRefresh := 100 * time.Millisecond

	b.Add(components.NewSeparatorWidget(0))
	b.Add(components.NewWlanWidget(1))

	b.Add(components.NewSeparatorWidget(2))
	b.Add(components.NewDiskUsageWidget(3))
	b.Add(components.NewMemoryWidget(4))
	b.Add(components.NewCpuWidget(5))

	b.Add(components.NewSeparatorWidget(6))
	b.Add(components.NewDateWidget(7))

	b.Add(components.NewSeparatorWidget(8))
	b.Add(components.NewPowerWidget(9, &powerRefresh))

	b.Add(components.NewSeparatorWidget(10))
	b.Add(components.NewVolumeWidget(11, &volumeRefresh))

	b.Add(components.NewSeparatorWidget(12))
	b.Add(components.NewLogoutWidget(13))

	for {
		m := <-b.Output
		fmt.Println(m + ",")
	}
}
