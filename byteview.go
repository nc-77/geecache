package geecache

type ByteView struct {
	b []byte
}

func (v ByteView) Len() int {
	return len(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

// ByteSlice returns a copy of the data as a byte slice.
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}
