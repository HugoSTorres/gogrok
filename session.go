package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// Session holds the data from a given session, along with some metadata.
type Session struct {
	// Data holds all the messages that have been recorded in the session.
	Data []Message

	handle *pcap.Handle
}

// Message is just a raw captured HTTP request or response. This allows us to
// use a slice of messages in the Session struct, which frees us from having
// to keep a count of the individual messages in the session (because users can
// just get it by calling len on the Data slice).
type Message []byte

// NewSession creates a capture session using pcap and returns it. If an error
// occurs creating the session, it is returned.
func NewSession(dev, filter string) (sess *Session, err error) {
	handle, err := pcap.OpenLive(dev, 1500, false, -1)
	if err != nil {
		return
	}

	err = handle.SetBPFFilter(filter)
	if err != nil {
		return
	}

	sess = &Session{Data: []Message{}, handle: handle}
	return
}

// Record opens up a gopacket packet source and iterates through the packets as they come in.
func (s *Session) Record(ch chan Message) error {
	src := gopacket.NewPacketSource(s.handle, s.handle.LinkType())

	for p := range src.Packets() {
		httpData := p.ApplicationLayer()
		if httpData == nil {
			continue
		}

		s.Data = append(s.Data, httpData.Payload())

		ch <- Message(httpData.Payload())
	}

	return nil
}
