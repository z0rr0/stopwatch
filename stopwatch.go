// Package main is console stopwatch tool.
//
// Control buttons:
// S - start/stop,
// Q     - quit,
// L - lap.
//
package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/nsf/termbox-go"
)

var (
	// Version is tool version
	Version = "0.0"
	// Revision is get revision
	Revision = "N/A"
	// Date is build date
	Date = "2017-01-01 00:00:00"
	// name is program name
	name = "Stopwatch"
)

func start(cutout <-chan bool) {
	startTime := time.Now()
	printTime := func(d time.Duration, format string) {
		fmt.Printf(format, d.Hours(), d.Minutes(), d.Seconds(), d.Nanoseconds())
	}
	ticker := time.NewTicker(20 * time.Millisecond) // 50Hz
	fmt.Println("--")
	for {
		select {
		case val, ok := <-cutout:
			printTime(time.Now().Sub(startTime), "\r%4.f:%02.f:%02.f.%v\n")
			if !ok || val {
				ticker.Stop()
				return
			}
		case <-ticker.C:
			printTime(time.Now().Sub(startTime), "\r%4.f:%02.f:%02.f.%v")
		}
	}
}

func watcher(control <-chan rune, ec chan<- error) {
	var isRun bool
	cutout := make(chan bool) // true is stop, false - lap
	defer close(cutout)
	for {
		select {
		case cmd, ok := <-control:
			if !ok {
				return
			}
			switch cmd {
			case 113, 81: // quit
				if isRun {
					cutout <- true
					isRun = false
				}
				ec <- nil
				return
			case 115, 83: // start/stop
				if isRun {
					cutout <- true
					isRun = false
				} else {
					go start(cutout)
					isRun = true
				}
			case 108, 76: // lap
				if isRun {
					cutout <- false
				}
			}
		}
	}
}

// listen listens keyboard eventsa and sends them to control channel.
func listen(control chan<- rune, ec chan<- error) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventError:
			ec <- ev.Err
		case termbox.EventKey:
			control <- ev.Ch
		}
	}
}

func main() {
	version := flag.Bool("version", false, "show version")
	flag.Parse()
	if *version {
		fmt.Printf("%v %v, %v, %v\n", name, Version, Revision, Date)
		return
	}

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputCurrent)

	erChan := make(chan error)
	controls := make(chan rune)
	defer close(erChan)
	defer close(controls)

	fmt.Println("S - start/stop, L - lap,  Q - quit")

	go listen(controls, erChan)
	go watcher(controls, erChan)

	err = <-erChan
	if err != nil {
		panic(err)
	}
}
