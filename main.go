package main

import (
	"log"
	"net"
	"net/http"

	"github.com/HugoSTorres/gogrok/pb"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	sess, err := NewSession()
	if err != nil {
		log.Fatalf("unable to open capture session: %s", err)
	}

	s := grpc.NewServer()
	pb.RegisterSnifferServer(s, sess)
	s.Serve(lis)
}

// func main() {
// 	log.Print("init gogrok")

// 	sess, err := NewSession()
// 	if err != nil {
// 		log.Fatalf("unable to open capture session: %s", err)
// 	}

// 	ch := make(chan Message)

// 	// Handle signals cleanly
// 	sig := make(chan os.Signal, 1)
// 	exit := make(chan struct{})
// 	signal.Notify(sig, os.Interrupt)

// 	log.Println("to quit gogrok, just press CTRL-C")

// 	go func() {
// 		for {
// 			select {
// 			case <-sig:
// 				exit <- struct{}{}
// 			}
// 		}
// 	}()

// 	// Let's give the sniffer three tries to start up, just in case something rare
// 	// and intermittent happens.
// 	go func() {
// 		for i := 0; i < 3; i++ {
// 			log.Fatalf("error recording: %s\n", sess.Record(ch))
// 		}
// 	}()

// 	http.HandleFunc("/", genStreamHandler(ch))

// 	// Let's give the server three tries to start up, just in case something rare
// 	// and intermittent happens.
// 	go func() {
// 		for i := 0; i < 3; i++ {
// 			log.Fatalf("error serving http: %s\n", http.ListenAndServe(":3000", nil))
// 		}
// 	}()

// 	<-exit

// 	log.Println("shutting down")

// 	close(exit)
// 	close(sig)
// 	close(ch)
// }

func genStreamHandler(ch <-chan Message) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("streamHandler - requesting capture data")

		m := <-ch
		w.Write([]byte(m))
		return
	}
}
