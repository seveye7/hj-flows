package flows

type Stream struct {
	c       chan [][]byte
	f       func(*Stream)
	reader  Reader // 源stream才有的数据
	gourpId string
	buffs   [][]byte
}

// SendToStream 发送数据到流，通过mq转发
func (s *Stream) SendToStream(topic string) {
}
