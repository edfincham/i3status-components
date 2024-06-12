package components

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	BATTERY_DISCHARGING = iota
	BATTERY_CHARGING
	BATTERY_FULL
)

type BatteryInfo struct {
	PercentRemaining float64
	SecondsRemaining float64
	status           int
}

type PowerWidget struct {
	BaseWidget
}

func NewPowerWidget(instance int, refresh *time.Duration) *PowerWidget {
	w := PowerWidget{
		BaseWidget: *NewBaseWidget(instance, refresh),
	}
	return &w
}

func (w *PowerWidget) getStatus() (string, string) {
	bi, _ := readBatteryInfo(0)
	remaining := ReadableDuration(int64(bi.SecondsRemaining))

	var prefix, colour string

	if bi.IsCharging() {
		colour = GREEN
		prefix = "\uf1e6"
		return fmt.Sprintf("%s %0.1f%%", prefix, bi.PercentRemaining), colour
	} else if bi.IsFull() {
		colour = GREEN
		prefix = "\uf240"
		return fmt.Sprintf("%s %0.1f%%", prefix, bi.PercentRemaining), colour
	} else {
		if bi.PercentRemaining < 10 {
			colour = RED
			prefix = "\uf244"
		} else if bi.PercentRemaining < 25 {
			colour = YELLOW
			prefix = "\uf243"
		} else if bi.PercentRemaining < 50 {
			colour = YELLOW
			prefix = "\uf243"
		} else if bi.PercentRemaining < 75 {
			colour = GREEN
			prefix = "\uf242"
		} else if bi.PercentRemaining <= 100 {
			colour = GREEN
			prefix = "\uf241"
		}

		return fmt.Sprintf("%s %0.1f%% (%s)", prefix, bi.PercentRemaining, remaining), colour
	}
}

func (w *PowerWidget) basicLoop() {
	msg := NewMessage()
	msg.Name = "power"
	msg.Colour = WHITE
	msg.Instance = strconv.Itoa(w.Instance)
	for {
		msg.FullText, msg.Colour = w.getStatus()
		w.Output <- *msg
		time.Sleep(w.Refresh)
	}
}

func (w *PowerWidget) Start() {
	go w.basicLoop()
	go w.readLoop()
}

func (batteryInfo *BatteryInfo) IsCharging() bool {
	return batteryInfo.status == BATTERY_CHARGING
}

func (batteryInfo *BatteryInfo) IsFull() bool {
	return batteryInfo.status == BATTERY_FULL
}

func readBatteryInfo(battery int) (*BatteryInfo, error) {
	rawInfo := make(map[string]string)
	batteryInfo := new(BatteryInfo)

	path := fmt.Sprintf("/sys/class/power_supply/BAT%d/uevent", battery)
	if !FileExists(path) {
		return batteryInfo, errors.New("Battery not found")
	}
	callback := func(line string) bool {
		data := strings.Split(string(line), "=")
		rawInfo[data[0]] = data[1]
		return true
	}
	ReadLines(path, callback)

	var remaining, present, full float64
	batteryInfo.status = BATTERY_DISCHARGING

	if rawInfo["POWER_SUPPLY_STATUS"] == "Charging" {
		batteryInfo.status = BATTERY_CHARGING
	} else if rawInfo["POWER_SUPPLY_STATUS"] == "Full" {
		batteryInfo.status = BATTERY_FULL
	}

	/* Convert to float shorthand */
	parseFloat := func(keys ...string) float64 {
		for _, key := range keys {
			if _, ok := rawInfo[key]; ok {
				f, _ := strconv.ParseFloat(rawInfo[key], 64)
				return f
			}
		}
		return 0.
	}

	/* Read values from file */
	remaining = parseFloat("POWER_SUPPLY_CHARGE_NOW")
	present = parseFloat("POWER_SUPPLY_CURRENT_NOW")
	full = parseFloat("POWER_SUPPLY_CHARGE_FULL")

	if full == 0 {
		return batteryInfo, errors.New("Battery full design missing")
	}

	batteryInfo.PercentRemaining = (remaining / full) * 100

	var remainingTime float64
	if present > 0 {
		if batteryInfo.status == BATTERY_DISCHARGING {
			remainingTime = remaining / present
			batteryInfo.SecondsRemaining = remainingTime * 3600
		}
	}
	return batteryInfo, nil
}
