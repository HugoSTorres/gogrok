package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
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

	go func() {
		err = sess.Record(ch)
		if err != nil {
			log.Fatalf("error recording: %s\n", err)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("requesting capture data")

		m := <-ch
		w.Write([]byte(m))
		return
	})

	// Let's give the server three tries to start up, just in case something rare
	// and intermittent happens.
	go func() {
		for i := 0; i < 3; i++ {
			log.Fatalf("error serving http: %s\n", http.ListenAndServe(":3000", nil))
		}
	}()

	<-exit

	log.Println("shutting down")

	close(exit)
	close(sig)
	close(ch)
}
