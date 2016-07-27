package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
)

const (
	port = ":50051"
)

func main() {
	log.Print("init gogrok")

	sess, err := NewSession()
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
