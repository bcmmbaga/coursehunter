package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

var stateFile = "state.json"

type state struct {
	Path   string  `json:"path"`
	Videos []video `json:"videos"`
}

func showProgress(file string, contentLen int64, done chan int64) {
	var (
		stop bool
		toMB = float64(1024 * 1024)
	)

	for {
		select {
		case <-done:
			stop = true
		default:
			f, err := os.OpenFile(file, os.O_RDONLY, 0555)
			if err != nil {
				log.Fatalf("Failed to open file for read: %v", err)
			}
			defer f.Close()

			info, err := f.Stat()
			if err != nil {
				log.Fatalf("Failed to retrieve file Info: %v", err)
			}

			size := info.Size()
			fmt.Printf("\033[2K\r Downloading %s -----> Received %.1fM out of %.1fM", path.Base(file),
				float64(size)/toMB, float64(contentLen)/toMB)
		}

		if stop {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func saveToDisk(reader io.Reader, writter io.Writer) int64 {
	written, err := io.Copy(writter, reader)
	if err != nil {
		log.Fatalf("Error while CopyToDisk: %v", err)

	}

	return written
}

func (h *hunter) saveState() {
	states := state{
		Path:   h.Path,
		Videos: h.Videos,
	}
	data, err := json.MarshalIndent(states, " ", " ")
	if err != nil {
		log.Fatalf("Failed to save state(s): %v", err)
	}

	ioutil.WriteFile(stateFile, data, 0600)
}
