package dispatcher

import (
	"errors"

	"github.com/bitly/go-nsq"
)

type Producer struct {
	nsqdEndpoints   []string
	nsqProducerList []*nsq.Producer
	index           int
	nsqConf         *nsq.Config
	logger          Logger
}

func NewProducer(nsqdEndpoints []string) *Producer {
	if len(nsqdEndpoints) == 0 {
		return nil
	}
	prod := &Producer{
		nsqdEndpoints:   nsqdEndpoints,
		nsqProducerList: make([]*nsq.Producer, len(nsqdEndpoints)),
		index:           0,
		nsqConf:         nsq.NewConfig(),
		logger:          &DefaultLogger{},
	}
	return prod
}

func (p *Producer) SetLogger(l Logger) {
	p.logger = l
}

func (p *Producer) Publish(e *Event) (err error) {
	if p == nil {
		return errors.New("producer is nil")
	}
	if e == nil {
		return errors.New("event is nil")
	}

	defer func() {
		if err != nil {
			p.index++
		}
	}()
	i := p.index % len(p.nsqProducerList)
	prod := p.nsqProducerList[i]
	if prod == nil {
		prod, err = nsq.NewProducer(p.nsqdEndpoints[i], p.nsqConf)
		if err != nil {
			return err
		}
		prod.SetLogger(nil, nsq.LogLevelError)
		p.nsqProducerList[i] = prod
	}
	bs, err := e.Marshal()
	if err != nil {
		return err
	}
	return prod.Publish(EventHeaderToString(e.Header), bs)
}
