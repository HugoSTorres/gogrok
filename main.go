package main

import (
	"flag"
	"fmt"
	"log"
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
		log.Fatalf("unable to open capture session: %s", err)
	}

	ch := make(chan Message)

	// Handle signals cleanly
	sig := make(chan os.Signal, 1)
	exit := make(chan struct{})
	signal.Notify(sig, os.Interrupt)

	log.Println("to quit gogrok, just press CTRL-C")

	go func() {
		for {
			select {
			case <-sig:
				exit <- struct{}{}
			}
		}
	}()

	// Let's give the sniffer three tries to start up, just in case something rare
	// and intermittent happens.
	go func() {
		for i := 0; i < 3; i++ {
			log.Fatalf("error recording: %s\n", sess.Record(ch))
		}
	}()

	go func() {
		for {
			msg := <-ch
			fmt.Println(string(msg))
		}
	}()

	<-exit

	log.Println("shutting down")

	close(exit)
	close(sig)
	close(ch)
}
