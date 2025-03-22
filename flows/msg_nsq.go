package flows

import (
	"context"
	"fmt"

	"github.com/nsqio/go-nsq"
)

type NsqConfig struct {
	Hosts    []string `yaml:"hosts"`
	BatchMax int64    `yaml:"batchMax"`
	Name     string
}

type NsqWriter struct {
	Config *NsqConfig
	writer *nsq.Producer
	c      chan *NsqMessage
	msgs   []NsqMessage
}

type NsqMessage struct {
	Topic string
	Value [][]byte
}

func NewNsqWriter(ctx context.Context, config *NsqConfig) Writer {
	kw := &NsqWriter{
		Config: config,
	}

	cfg := nsq.NewConfig()
	cfg.UserAgent = fmt.Sprintf("%v go-nsq/%s", config.Name, nsq.VERSION)

	producer, err := nsq.NewProducer(config.Hosts[0], cfg)
	if err != nil {
		return nil
	}
	kw.writer = producer

	go func() {
		for {
			select {
			case <-ctx.Done():
				kw.write(true)
				kw.writer.Stop()
				return
			case m := <-kw.c:
				kw.msgs = append(kw.msgs, *m)
				kw.write(false)
			}
		}
	}()

	return kw
}

func (w *NsqWriter) write(force bool) {
	if !force && len(w.msgs) < int(w.Config.BatchMax) {
		return
	}
	if len(w.msgs) > 0 {
		for _, v := range w.msgs {
			w.writer.MultiPublish(v.Topic, v.Value)
		}
		w.msgs = w.msgs[:0]
	}
}

func (w *NsqWriter) Write(ctx context.Context, topic string, bs [][]byte) {
	w.c <- &NsqMessage{
		Topic: topic,
		Value: bs,
	}
}

type NsqReader struct {
	Config *NsqConfig
	topic  string
	reader *nsq.Consumer
	c      chan []byte
}

func NewNsqReader(config *NsqConfig, selfGourpId, topic string) Reader {
	nsqConfig := nsq.NewConfig()
	nsqConfig.MaxInFlight = 100
	consumer, err := nsq.NewConsumer(topic, selfGourpId, nsqConfig)
	if err != nil {
		panic(err)
	}

	err = consumer.ConnectToNSQDs(config.Hosts)

	nr := &NsqReader{
		c:      make(chan []byte),
		Config: config,
		topic:  topic,
		reader: consumer,
	}

	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		nr.c <- message.Body
		return nil
	}))
	if err != nil {
		panic(err)
	}
	return nr
}

func (r *NsqReader) Read(ctx context.Context) ([][]byte, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case buff := <-r.c:
			return [][]byte{buff}, nil
		}
	}
}
