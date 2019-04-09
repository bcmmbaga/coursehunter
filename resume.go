package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func getStates() *state {
	f, err := os.OpenFile(stateFile, os.O_RDONLY, 0600)
	if err != nil {
		log.Fatalf("Error while opening file: %v", err)
	}
	defer f.Close()

	states := &state{}
	err = json.NewDecoder(f).Decode(states)
	if err != nil {
		log.Fatalf("Error while decoding state(s): %v", err)
	}

	return states
}

func (s *state) resume(start int, notify chan os.Signal) {
	for i := start - 1; i < len(s.Videos); i++ {
		go func() {
			<-notify
			os.Exit(1)
		}()

		downloadRange(fmt.Sprintf("%s/%s.mp4", s.Path, s.Videos[i].Name), s.Videos[i].URL)
	}
}

func downloadRange(file, url string) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalf("Failed to create new request: %v", err)
	}

	size := fileSize(file)
	ranges := fmt.Sprintf("bytes=%d-", size)
	req.Header.Add("Range", ranges)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error while resuming video: %v", err)
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		log.Fatalf("Failed to open file for read: %v", err)
	}
	defer f.Close()

	done := make(chan int64)

	go showProgress(file, resp.ContentLength+size, done)
	done <- saveToDisk(resp.Body, f)
}

func fileSize(file string) int64 {
	f, err := os.OpenFile(file, os.O_RDONLY, 0555)
	if err != nil {
		return 0
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		log.Fatalf("Failed to retrieve file Info: %v", err)
	}

	return info.Size()
}
