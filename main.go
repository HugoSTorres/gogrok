package main

import "log"

func main() {
	log.Print("init gogrok")

	sess, err := NewSession()
	if err != nil {
		log.Fatalf("unable to open capture session: %s", err)
	}

	log.Fatalf("error recording: %s\n", sess.Record())
}
