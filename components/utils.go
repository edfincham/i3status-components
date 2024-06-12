// this has been stolen form
// https://github.com/coleifer/mastodon/blob/master/utils.go
package components

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	RED        = "#ff0000"
	YELLOW     = "#ffff00"
	GREEN      = "#43b33b"
	WHITE      = "#ffffff"
	BLUE       = "#0000ff"
	BACKGROUND = "#363537"
)

var REFRESH = 1000 * time.Millisecond

func ReadLines(fileName string, callback func(string) bool) {
	fin, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "The file %s does not exist!\n", fileName)
		return
	}
	defer fin.Close()

	reader := bufio.NewReader(fin)
	for line, _, err := reader.ReadLine(); err != io.EOF; line, _, err = reader.ReadLine() {
		if !callback(string(line)) {
			break
		}
	}
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func MakeBar(percent float64, bar_size int) string {
	var bar bytes.Buffer
	cutoff := int(percent * .01 * float64(bar_size))
	bar.WriteString("[")
	for i := 0; i < bar_size; i += 1 {
		if i <= cutoff {
			bar.WriteString("#")
		} else {
			bar.WriteString(" ")
		}
	}
	bar.WriteString("]")
	return bar.String()
}

func ReadableDuration(n int64) string {
	hours := n / 3600
	minutes := (n % 3600) / 60
	return fmt.Sprintf("%d:%02d", hours, minutes)
}
