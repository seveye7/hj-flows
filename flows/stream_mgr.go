package flows

import (
	"context"
	"log"

	"hj-flows/utils"

	"golang.org/x/sync/errgroup"
)

type StreamMgr struct {
	ctx         context.Context
	streams     map[any]*Stream
	wg          errgroup.Group
	cancel      func()
	kafkaConfig *KafkaConfig
	nsqConfig   *NsqConfig
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

	return m
}

// option 流处理参数
type Option func(*StreamMgr)

func WithKafka(config *KafkaConfig) Option {
	return func(m *StreamMgr) {
		m.kafkaConfig = config
		m.writer = NewKafkaWriter(m.ctx, m.kafkaConfig)
	}
}

func WithNsq(config *NsqConfig) Option {
	return func(m *StreamMgr) {
		m.nsqConfig = config
		m.writer = NewNsqWriter(m.ctx, m.nsqConfig)
	}
}

// WithTopic 注册topic和流处理函数
func WithTopic(gourpId, topic string, f func(*Stream)) Option {
	return func(m *StreamMgr) {
		// TODO: 支持同时注册一个topic，groupI不一样
		if _, ok := m.streams[topic]; ok {
			panic("topic already registered: " + topic)
		}
		log.Println("WithTopic", topic)
		s := &Stream{f: f, c: make(chan [][]byte, 128), streamMgr: m, gourpId: gourpId}
		if m.kafkaConfig != nil {
			s.reader = NewKafkaReader(m.kafkaConfig, gourpId, topic)
		} else if m.nsqConfig != nil {
			s.reader = NewNsqReader(m.nsqConfig, gourpId, topic)
		}
		m.streams[topic] = s
	}
}

func (m *StreamMgr) GetStream(topic string) *Stream {
	return m.streams[topic]
}

// Start 启动所有的流
func (m *StreamMgr) Start() {
	for _, stream := range m.streams {
		m.wg.Go(func() error {
			if utils.IsNil(stream.reader) {
				return nil
			}
			for {
				buffs, err := stream.reader.Read(m.ctx)
				if err != nil {
					break
				}
				stream.c <- buffs
			}
			return nil
		})
		m.wg.Go(func() error {
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
