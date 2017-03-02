package main

import (
	"bytes"
	"flag"
	"fmt"
	"git.meiqia.com/infrastructure/event-dispatcher"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	lookupdHttpAddress string
	lookupdTcpAddress  string
	header             string
)

func main() {
	flag.StringVar(&lookupdHttpAddress, "lookup-http-address", "127.0.0.1:4161", "lookupd http address")
	flag.StringVar(&lookupdTcpAddress, "lookup-tcp-address", "127.0.0.1:4160", "lookupd tcp address")
	flag.StringVar(&header, "header", "king-noble-knight-peasant", "use slash to seperate")
	flag.Parse()
	eh := stringToEventHeader(header)

	cons := dispatcher.NewConsumer([]string{lookupdTcpAddress}, []string{lookupdHttpAddress}, *eh)
	cons.SetHandler(handleMessage)
	cons.Start()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

func handleMessage(e *dispatcher.Event) {
	log.Printf("Received from [%s] topic, payload is [%s]\n", eventHeaderToString(e.Header), string(e.Body.Content))
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

func eventHeaderToString(e *dispatcher.EventHeader) string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s-%s-%s-%s", e.King, e.Noble, e.Knight, e.Peasant))
	for _, tag := range e.Tags {
		buffer.WriteString(fmt.Sprintf("-%s", tag))
	}

	return buffer.String()
}
