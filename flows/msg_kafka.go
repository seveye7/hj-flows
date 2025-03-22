package flows

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type Writer interface {
	Write(context.Context, string, [][]byte)
}
type Reader interface {
	Read(context.Context) ([][]byte, error)
}

type KafkaConfig struct {
	Hosts    []string `yaml:"hosts"`
	BatchMax int64    `yaml:"batchMax"`
}

type KafkaWriter struct {
	Config *KafkaConfig
	writer *kafka.Writer
	c      chan *kafka.Message
	msgs   []kafka.Message
}

func NewKafkaWriter(ctx context.Context, config *KafkaConfig) Writer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(config.Hosts...),
		AllowAutoTopicCreation: true,
		// RequiredAcks:           kafka.RequireAll,
	}

	kw := &KafkaWriter{
		Config: config,
		writer: w,
		c:      make(chan *kafka.Message, 1024),
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				kw.write(true)
				w.Close()
				return
			case m := <-kw.c:
				kw.msgs = append(kw.msgs, *m)
				kw.write(false)
			}
		}
	}()

	return kw
}

func (w *KafkaWriter) write(force bool) {
	if !force && len(w.msgs) < int(w.Config.BatchMax) {
		return
	}
	if len(w.msgs) > 0 {
		w.writer.WriteMessages(context.Background(), w.msgs...)
		w.msgs = w.msgs[:0]
	}
}

func (w *KafkaWriter) Write(ctx context.Context, topic string, bs [][]byte) {
	for _, v := range bs {
		w.c <- &kafka.Message{
			Topic: topic,
			Value: v,
		}
	}
}

type KafkaReader struct {
	Config *KafkaConfig
	topic  string
	reader *kafka.Reader
}

func NewKafkaReader(config *KafkaConfig, selfGourpId, topic string) Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: config.Hosts,
		GroupID: selfGourpId,
		Topic:   topic,
		Dialer: &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
		},
		MaxBytes: 10e6, // 10MB
	})
	return &KafkaReader{
		Config: config,
		topic:  topic,
		reader: r,
	}
}

func (r *KafkaReader) Read(ctx context.Context) ([][]byte, error) {
	m, err := r.reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	return [][]byte{m.Value}, nil
}
