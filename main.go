package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
)

const use = "usage: gogrok --interface <device> <filter>\n\n"

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, use)
		flag.PrintDefaults()
	}
}

func main() {
	var dev string
	flag.StringVar(&dev, "i", "", "the device to listen on")

	flag.Parse()

	if len(dev) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	filter := "port 80"
	if len(flag.Args()) != 0 {
		filter = strings.Join(flag.Args(), " ")
	}

	sess, err := NewSession(dev, filter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to open capture session: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("gogrok: listening on %s for %s\n", dev, filter)

	ch := make(chan Message)

	// Handle signals cleanly
	sig := make(chan os.Signal, 1)
	exit := make(chan struct{})
	signal.Notify(sig, os.Interrupt)

	fmt.Println("gogrok: CTRL-C to quit")

	go func() {
		for {
			select {
			case <-sig:
				fmt.Println("\ngogrok: shutting down")
				exit <- struct{}{}
			}
		}
	}()

	// Let's give the sniffer three tries to start up, just in case something rare
	// and intermittent happens.
	go func() {
		for i := 0; i < 2; i++ {
			fmt.Printf("error recording: %v\n", sess.Record(ch))
			fmt.Println("retrying...")
		}

		fmt.Fprintf(os.Stderr, "error recording: %s\n", sess.Record(ch))
	}()

	go func() {
		for {
			msg := <-ch
			fmt.Println(string(msg))
		}
	}()

	<-exit
}
