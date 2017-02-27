package dispatcher

import (
	"github.com/bitly/go-nsq"
	"github.com/golang/glog"
)

type EventHandler func(e *Event)

type Consumer struct {
	conf   *Config
	cons   []*nsq.Consumer
	ehs    EventHandler
	header EventHeader
}

func NewConsumer(conf *Config, king, noble, knight, peasant string) *Consumer {
	c := &Consumer{
		conf:   conf,
		header: EventHeader{king, noble, knight, peasant},
	}

	return c
}

func (c *Consumer) newNsqConsumer(endpoints []string, topic, channel string) *nsq.Consumer {
	cons, err := nsq.NewConsumer(topic, channel, nsq.NewConfig())
	if err != nil {
		glog.Error("NewConsumer failed: ", err)
	}

	cons.AddHandler(c)

	err = cons.ConnectToNSQLookupds(endpoints)
	if err != nil {
		glog.Error("Connect to nsq lookup error: ", err)
	}
	cons.AddHandler(c)

	return cons
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
	//TODO:
}
