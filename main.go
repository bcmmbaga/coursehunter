package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	course   string
	email    string
	password string
	start    int
)

func init() {
	flag.StringVar(&course, "n", "", "coursename to download")
	flag.StringVar(&email, "e", "", "coursehunter.net login email")
	flag.StringVar(&password, "p", "", "coursehunter.net login password")
	flag.IntVar(&start, "start", 1, "index video to resume from")
	flag.Parse()

}
func main() {
	var (
		h   *hunter
		err error
	)

	notify := make(chan os.Signal, 1)
	signal.Notify(notify, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	args := flag.Args()
	if len(args) == 0 {
		h, err = newHunter(course, email, password)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		mkdirCourse(course)
		h.download(notify)

	} else if len(args) >= 2 && args[0] == "resume" {
		getStates().resume(start, notify)

	} else {
		printUsage()
		os.Exit(1)

	}
}

func mkdirCourse(course string) {
	if _, err := os.Stat(course); os.IsNotExist(err) {
		err := os.Mkdir(course, 0755)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func printUsage() {
	helpstr := `
[INFO]: coursehunter
	
[USAGE]: coursehunter [command] [options...] 
	COMMAND:
		resume resume interrupted downloads
	OPTIONS:
		-n coursename
		-e email
		-p password
		- start   
`
	fmt.Println(helpstr)
}
