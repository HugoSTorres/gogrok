package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
)

const use = "usage: gogrok [--src_addr=<ip>] [--src_port=<port>] [--dest_addr=<ip>] [--dest_port=<port>]"

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, use)
		flag.PrintDefaults()
	}
}

func main() {
	var dev, srcAddr, destAddr string
	var srcPort, destPort int
	flag.StringVar(&dev, "device", "", "the device to listen on")
	flag.StringVar(&srcAddr, "src_addr", "", "the client's ip address (without port)")
	flag.StringVar(&destAddr, "dest_addr", "", "the server's ip address (without port)")
	flag.IntVar(&srcPort, "src_port", 80, "the client's port")
	flag.IntVar(&destPort, "dest_port", -1, "the server's port")

	flag.Parse()

	if len(dev) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	sess, err := NewSession(dev, srcAddr, destAddr, srcPort, destPort)
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
