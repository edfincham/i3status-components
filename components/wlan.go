package components

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type WlanInfo struct {
	PublicIP      string
	NetworkStatus string
	NetworkName   string
	color         string
}

type WlanWidget struct {
	BaseWidget
}

func NewWlanWidget(instance int) *WlanWidget {
	w := WlanWidget{
		BaseWidget: *NewBaseWidget(instance, nil),
	}
	return &w
}

func (w *WlanWidget) basicLoop() {
	msg := NewMessage()
	msg.Name = "wlan"
	msg.Colour = WHITE
	msg.Instance = strconv.Itoa(w.Instance)
	for {
		msg.FullText, msg.Colour = w.getStatus()
		w.Output <- *msg
		time.Sleep(5000 * time.Millisecond)
	}
}

func (w *WlanWidget) getStatus() (string, string) {
	wi, _ := getIP()

	var prefix, colour string
	colour = WHITE
	prefix = "\uf1eb"

	return fmt.Sprintf("%s %s", prefix, wi.PublicIP), colour
}

func (w *WlanWidget) Start() {
	go w.basicLoop()
	go w.readLoop()
}

func getIP() (*WlanInfo, error) {
	wlanInfo := new(WlanInfo)

	url := "https://ifconfig.co"
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error: %v", err)
		return wlanInfo, nil
	}
	defer resp.Body.Close()

	ipAddress, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error: %v", err)
		return wlanInfo, nil
	}

	wlanInfo.PublicIP = strings.Trim(string(ipAddress), "\n")
	return wlanInfo, nil
}
