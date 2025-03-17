package flows

import (
	"context"

	"github.com/go-pay/errgroup"
)

type StreamMgr struct {
	ctx         context.Context
	streams     map[any]*Stream
	wg          errgroup.Group
	cancel      func()
	kafkaConfig *KafkaConfig
	writer      Writer
}

func NewStreamMgr(options ...Option) *StreamMgr {
	ctx, cancel := context.WithCancel(context.Background())
	m := &StreamMgr{
		streams: make(map[any]*Stream),
		ctx:     ctx,
		cancel:  cancel,
	}
	for _, option := range options {
		option(m)
	}

	if m.kafkaConfig != nil {
		m.writer = NewKafkaWriter(ctx, m.kafkaConfig)
	}

	return m
}

// option 流处理参数
type Option func(*StreamMgr)

func WithKafka(config *KafkaConfig) Option {
	return func(m *StreamMgr) {
		m.kafkaConfig = config
	}
}

// WithTopic 注册topic和流处理函数
func WithTopic(gourpId, topic string, f func(*Stream)) Option {
	return func(m *StreamMgr) {
		if _, ok := m.streams[topic]; ok {
			panic("topic already registered: " + topic)
		}
		s := &Stream{f: f, c: make(chan [][]byte, 128)}
		if m.kafkaConfig != nil {
			s.gourpId = gourpId
			s.reader = NewKafkaReader(m.kafkaConfig, gourpId, topic)
		}
		m.streams[topic] = s
	}
}

// Start 启动所有的流
func (m *StreamMgr) Start() {
	for _, stream := range m.streams {
		m.wg.Go(func(ctx context.Context) error {
			for {
				buffs, err := stream.reader.Read(m.ctx)
				if err != nil {
					break
				}
				stream.c <- buffs
			}
			return nil
		})
		m.wg.Go(func(ctx context.Context) error {
			for {
				select {
				case <-m.ctx.Done():
					return nil
				case buffs := <-stream.c:
					stream.buffs = buffs
					stream.f(stream)
				}
			}
		})
	}
}

// Stop 停止所有的流
func (m *StreamMgr) Stop() {
	m.cancel()
	m.wg.Wait()
}
