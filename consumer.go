package dispatcher

import (
	"github.com/bitly/go-nsq"
	"time"
)

type EventHandler func(e *Event)

type Consumer struct {
	lookupTcpEndpoints  []string
	lookupHttpEndpoints []string
	consumers           []*nsq.Consumer
	interestedTopics    map[string]bool
	ehs                 EventHandler
	interested          EventHeader
	logger              Logger
	nsqChannel          string
}

func NewConsumer(lookupTcpEndpoints, lookupHttpEndpoints []string, h EventHeader) *Consumer {
	return NewConsumerWithChannel(lookupTcpEndpoints, lookupHttpEndpoints, h, defaultChannel)
}

func NewConsumerWithChannel(lookupTcpEndpoints, lookupHttpEndpoints []string, h EventHeader, ch string) *Consumer {
	c := &Consumer{
		lookupTcpEndpoints:  lookupTcpEndpoints,
		lookupHttpEndpoints: lookupHttpEndpoints,
		interestedTopics:    make(map[string]bool),
		interested:          h,
		consumers:           []*nsq.Consumer{},
		logger:              &DefaultLogger{},
		nsqChannel:          ch,
	}

	return c
}

func (c *Consumer) SetLogger(l Logger) {
	c.logger = l
}

func (c *Consumer) newNsqConsumer(endpoints []string, topic, channel string) {
	cons, err := nsq.NewConsumer(topic, channel, nsq.NewConfig())
	if err != nil {
		c.logger.Error("NewConsumer failed: ", err)
		return
	}

	cons.AddHandler(c)
	cons.SetLogger(nil, nsq.LogLevelError)

	err = cons.ConnectToNSQLookupds(endpoints)
	if err != nil {
		c.logger.Error("Connect to nsq lookup error: ", err)
		return
	}

	c.consumers = append(c.consumers, cons)
	c.interestedTopics[topic] = true
}

func (c *Consumer) SetHandler(eh EventHandler) {
	c.ehs = eh
}

func (c *Consumer) Start() {
	go c.listen()
}

func (c *Consumer) Stop() {
	//TODO
}

func (c *Consumer) HandleMessage(msg *nsq.Message) error {
	event := &Event{}
	err := event.Unmarshal(msg.Body)
	if err != nil {
		return err
	}
	if c.ehs != nil {
		c.ehs(event)
	}

	return nil
}

func (c *Consumer) listen() {
	for {
		if topics, err := listTopics(c.lookupHttpEndpoints); err == nil {
			for _, topic := range topics {
				if _, ok := c.interestedTopics[topic]; ok {
					continue
				}
				if !c.shouldInterest(topic) {
					continue
				}
				c.logger.Debug("newly interested topic", topic)
				c.newNsqConsumer(c.lookupHttpEndpoints, topic, c.nsqChannel)
			}
		} else {
			c.logger.Error(err)
		}
		time.Sleep(updateTopicsDur)
	}
}

/*
Strategy:
	1. interested King is empty string match all the case
*/
func (c *Consumer) shouldInterest(topic string) bool {
	eventHeader, err := StringToEventHeader(topic)
	if err != nil {
		return false
	}

	if determined, result := determine(c.interested.King, eventHeader.King); determined {
		return result
	}
	if determined, result := determine(c.interested.Noble, eventHeader.Noble); determined {
		return result
	}
	if determined, result := determine(c.interested.Knight, eventHeader.Knight); determined {
		return result
	}
	if determined, result := determine(c.interested.Peasant, eventHeader.Peasant); determined {
		return result
	}

	if len(c.interested.Tags) > len(eventHeader.Tags) {
		return false
	}
	for i, tag := range c.interested.Tags {
		if determined, result := determine(tag, eventHeader.Tags[i]); determined {
			return result
		}
	}

	return true
}

func determine(src string, dest string) (determinable, result bool) {
	if len(src) == 0 {
		return true, true
	}
	if src != dest {
		return true, false
	}
	return false, true
}
