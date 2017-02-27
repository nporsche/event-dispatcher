package dispatcher

import (
	"errors"

	"github.com/bitly/go-nsq"
	"github.com/golang/glog"
)

type Producer struct {
	conf            *Config
	nsqProducerList []*nsq.Producer
	index           int
	nsqConf         *nsq.Config
}

func NewProducer(conf *Config) *Producer {
	if conf == nil {
		glog.Error("config is nil")
		return nil
	}
	if conf.Endpoints == nil || len(conf.Endpoints) == 0 {
		glog.Error("invalid endpoints in config")
		return nil
	}
	prod := &Producer{
		conf:            conf,
		nsqProducerList: make([]*nsq.Producer, len(conf.Endpoints)),
		index:           0,
		nsqConf:         nsq.NewConfig(),
	}
	return prod
}

func (p *Producer) Publish(e *Event) (err error) {
	if p == nil {
		glog.Error("producer is nil")
		return errors.New("producer is nil")
	}
	if e == nil {
		glog.Error("event is nil")
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
		prod, err = nsq.NewProducer(p.conf.Endpoints[i], p.nsqConf)
		if err != nil {
			return err
		}
		p.nsqProducerList[i] = prod
	}
	var bs []byte
	bs, err = e.Marshal()
	if err != nil {
		return err
	}

	err = prod.Publish(DefaultTopicNameBuilder(*e.Header), bs)
	return
}
