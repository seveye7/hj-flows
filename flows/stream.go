package flows

import (
	"context"

	"hj-flows/utils"
)

type Stream struct {
	c         chan [][]byte
	f         func(*Stream)
	reader    Reader // 源stream才有的数据
	streamMgr *StreamMgr
	gourpId   string
	buffs     [][]byte
}

// SendToStream 发送数据到流，通过mq转发
func (s *Stream) SendToStream(topic string) {
	if utils.IsNil(s.streamMgr.writer) {
		s.streamMgr.GetStream(topic).sendBuffs(s.buffs...)
		return
	}
	s.streamMgr.writer.Write(context.Background(), topic, s.buffs)
}

func (s *Stream) SendMessage(a any) {
	buff := Marshal(a)
	s.sendBuffs(utils.S2b(buff[0]))
}

func (s *Stream) sendBuffs(buff ...[]byte) {
	if s == nil {
		return
	}
	s.c <- buff
}
