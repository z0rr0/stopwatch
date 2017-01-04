// Package main is console stopwatch tool.
//
// Control buttons:
// S - start/stop,
// Q     - quit,
// L - lap.
//
package main

import (
	"errors"
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
	name = "Stopwatch"

	errStop = errors.New("stop")
)

func start(cutout <-chan bool) {
	startTime := time.Now()
	ticker := time.NewTicker(10 * time.Millisecond)
	fmt.Println("\nstart")
	for {
		select {
		case val := <-cutout:
			if val {
				ticker.Stop()
				return
			}
			d := time.Now().Sub(startTime)
			fmt.Printf("\r%4.f:%02.f:%02.f.%v\n", d.Hours(), d.Minutes(), d.Seconds(), d.Nanoseconds())
		case <-ticker.C:
			d := time.Now().Sub(startTime)
			fmt.Printf("\r%4.f:%02.f:%02.f.%v", d.Hours(), d.Minutes(), d.Seconds(), d.Nanoseconds())
		}
	}
}

func watcher(control <-chan rune, ec chan<- error) {
	var isRun bool
	cutout := make(chan bool) // true is stop, false - lap
	for {
		select {
		case cmd := <-control:
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
	termbox.SetInputMode(termbox.InputEsc)

	erch := make(chan error)
	controls := make(chan rune)

	fmt.Println("S - start/stop, L - lap,  Q - quit")

	go listen(controls, erch)
	go watcher(controls, erch)

	err = <-erch
	if err != nil {
		panic(err)
	}
}
