package stream

import (
	"bytes"
	"encoding/binary"
	"io"
)

func New(data []byte) *Stream {
	return &Stream{
		Reader: bytes.NewReader(data),
		order:  binary.LittleEndian,
	}
}

type Stream struct {
	order binary.ByteOrder
	*bytes.Reader
}

func (s *Stream) PeekChar() byte {
	c, _ := s.ReadByte()
	_ = s.UnreadByte()

	return c
}

func (s *Stream) ReadFString(n int) string {
	data := make([]byte, n)

	_, _ = s.Read(data)

	for i := range data {
		if data[i] == 0 {
			return string(data[:i])
		}
	}

	return string(data)
}

func (s *Stream) ReadCString() string {
	start := s.Pos()

	for {
		c, err := s.ReadByte()
		if c == 0 || err != nil {
			break
		}
	}

	strLen := s.Pos() - start

	_, _ = s.Seek(start, io.SeekStart)

	data := make([]byte, strLen)
	n, _ := s.Read(data)

	data = data[:n]

	if len(data) > 0 && data[len(data)-1] == 0 {
		data = data[:len(data)-1]
	}

	return string(data)
}

func (s *Stream) Pos() int64 {
	pos, _ := s.Seek(0, io.SeekCurrent)
	return pos
}

func (s *Stream) ReadInt16() int {
	var v int16

	_ = binary.Read(s, s.order, &v)

	return int(v)
}

func (s *Stream) ReadInt32() int {
	var v int32

	_ = binary.Read(s, s.order, &v)

	return int(v)
}

func (s *Stream) ReadObject(v any) error {
	return binary.Read(s, s.order, v)
}

func (s *Stream) SetOrder(order binary.ByteOrder) {
	s.order = order
}
