package main

import (
	"log"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func main() {
	log.Print("init gogrok")

	handle, err := pcap.OpenLive("en0", 1500, false, -1)
	if err != nil {
		log.Fatalf("error opening handle: %s", err)
	}

	err = handle.SetBPFFilter("tcp src port 80")
	if err != nil {
		log.Fatalf("error setting filter: %s", err)
	}

	f, err := os.Create("out.txt")
	if err != nil {
		log.Fatalf("error creating file for packet dump: %s", err)
	}

	src := gopacket.NewPacketSource(handle, handle.LinkType())
	for p := range src.Packets() {
		log.Printf("recv packet: %v", len(p.Data()))

		_, err := f.Write(p.Data())
		if err != nil {
			log.Fatalf("error writing to output file: %s", err)
		}
	}
}
