package main

import (
	"log"
	"net/http"
)

func main() {
	log.Print("init gogrok")

	sess, err := NewSession()
	if err != nil {
		log.Fatalf("unable to open capture session: %s", err)
	}

	ch := make(chan Message)

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

	log.Fatalf("error serving http: %s\n", http.ListenAndServe(":3000", nil))
}
