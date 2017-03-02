package main

import (
	"flag"
	"git.meiqia.com/infrastructure/event-dispatcher"
	"log"
	"strings"
)

var (
	nsqd   string
	header string
	body   string
)

func main() {
	flag.StringVar(&nsqd, "nsqd-address", "127.0.0.1:4150", "nsq address")
	flag.StringVar(&header, "event-header", "king-noble-knight-peasant-tag1-tag2", "event header")
	flag.StringVar(&body, "event-body", "test body message", "event body")
	flag.Parse()
	prod := dispatcher.NewProducer([]string{nsqd})
	event := &dispatcher.Event{
		Header: stringToEventHeader(header),
		Body:   &dispatcher.EventBody{Content: []byte(body)},
	}
	prod.Publish(event)
	log.Println("publist a message")
}

func stringToEventHeader(s string) *dispatcher.EventHeader {
	strs := strings.Split(s, "-")
	eh := &dispatcher.EventHeader{
		King:    strs[0],
		Noble:   strs[1],
		Knight:  strs[2],
		Peasant: strs[3],
		Tags:    strs[4:],
	}
	return eh
}
