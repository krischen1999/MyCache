package MyCache

/*缓存值的抽象与封装*/
type ByteView struct {
	data []byte
	//data 只读
	//data 存储缓存值
}

func (b ByteView) Len() int {
	return len(b.data)
}

func (b ByteView) ByteSlice() []byte {
	return CloneBytes(b.data)
}

func (b ByteView) String() string {
	return string(b.data)
}

func CloneBytes(data []byte) []byte {
	a := make([]byte, len(data))
	copy(a, data)
	return a
}
