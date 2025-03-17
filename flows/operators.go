package flows

// Filter 输入一个元素同时输出一个bool值为每个元素执行一个布尔 function，并保留那些 function 输出值为 true 的元素。
func Filter[T any](s *Stream, f func(*T) bool) *Stream {
	newStream := &Stream{}
	data, _, _ := UnMarshalBytes[T](s.buffs)
	for _, v := range data {
		if f(v) {
			newStream.buffs = append(newStream.buffs, MarshalBytes(v)...)
		}
	}
	return newStream
}

// Map 输入一个元素同时输出一个元素
func Map[T1 any, T2 any](s *Stream, f func(*T1) *T2) *Stream {
	newStream := &Stream{}
	data, _, _ := UnMarshalBytes[T1](s.buffs)
	for _, v := range data {
		newStream.buffs = append(newStream.buffs, MarshalBytes(f(v))...)
	}
	return newStream
}

// FlatMap 输入一个元素同时产生零个、一个或多个元素
func FlatMap[T1 any, T2 any](s *Stream, f func(*T1) []*T2) *Stream {
	newStream := &Stream{}
	data, _, _ := UnMarshalBytes[T1](s.buffs)
	for _, v := range data {
		newStream.buffs = append(newStream.buffs, MarshalBytes(f(v))...)
	}
	return newStream
}

// KeyBy 在逻辑上将流划分为不相交的分区。具有相同 key 的记录都分配到同一个分区。在内部， keyBy() 是通过哈希分区实现的。有多种指定 key 的方式。
func KeyBy[T any](s *Stream, f func(*T) int, n int) []*Stream {
	newstems := make([]*Stream, n)
	for i := 0; i < n; i++ {
		newstems[i] = &Stream{}
	}
	data, _, _ := UnMarshalBytes[T](s.buffs)
	for _, v := range data {
		i := f(v)
		newstems[i].buffs = append(newstems[i].buffs, MarshalBytes(v)...)
	}
	return newstems
}

// Reduce 在相同 key 的数据流上“滚动”执行 reduce。将当前元素与最后一次 reduce 得到的值组合然后输出新值。
func Reduce[T any](s *Stream, f func(*T, *T) *T) *Stream {
	newStream := &Stream{}
	data, _, _ := UnMarshalBytes[T](s.buffs)
	result := data[0]
	for i := 1; i < len(data); i++ {
		result = f(result, data[i])
	}
	newStream.buffs = append(newStream.buffs, MarshalBytes(result)...)
	return newStream
}

// Union 将两个或多个数据流联合来创建一个包含所有流中数据的新流。注意：如果一个数据流和自身进行联合，这个流中的每个数据将在合并后的流中出现两次。
func Union[T any](ss ...*Stream) *Stream {
	newStream := &Stream{}
	for _, v := range ss {
		newStream.buffs = append(newStream.buffs, v.buffs...)
	}
	return newStream
}
