package xmlb

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/GlebYaltchik/sc-keybind-extract/internal/stream"
)

type Node struct {
	NodeNameOffset      int32
	ContentOffset       int32
	AttributeCount      int16
	ChildCount          int16
	ParentNodeID        int32
	FirstAttributeIndex int32
	FirstChildIndex     int32
	Reserved            int32
}

type Reference struct {
	NameOffset  int32
	ValueOffset int32
}

type DataMap map[int32]string

type info struct {
	NodeTableOffset      int32
	NodeTableCount       int32
	AttributeTableOffset int32
	AttributeTableCount  int32
	ChildTableOffset     int32
	ChildTableCount      int32
	StringTableOffset    int32
	StringTableCount     int32
}

type Parser struct {
	s     *stream.Stream
	info  info
	names DataMap
	attrs []Reference
}

func NewParser(s *stream.Stream) (*Parser, error) {
	if s.PeekChar() != 'C' {
		return nil, fmt.Errorf("unknown data format")
	}

	hdr := s.ReadFString(7)

	switch hdr {
	case "CryXmlB", "CryXml":
		_ = s.ReadCString()
	case "CRY3SDK":
		_, _ = s.ReadByte()
		_, _ = s.ReadByte()
	default:
		return nil, fmt.Errorf("unknown data format")
	}

	headerPos := s.Pos()

	fileLength := s.ReadInt32()
	if int64(fileLength) != s.Size() {
		s.SetOrder(binary.LittleEndian)
		_, _ = s.Seek(headerPos, io.SeekStart)
		fileLength = s.ReadInt32()
	}

	if int64(fileLength) != s.Size() {
		return nil, fmt.Errorf("file length mismatch")
	}

	info := mustRead[info](s)

	return &Parser{
		s:     s,
		info:  info,
		names: readDict(s, int64(info.StringTableOffset)),
		attrs: readAttrs(s, int64(info.AttributeTableOffset), int(info.AttributeTableCount)),
	}, nil
}

func (p *Parser) ReadNodes() []Node {
	_, _ = p.s.Seek(int64(p.info.NodeTableOffset), io.SeekStart)

	nodes := make([]Node, p.info.NodeTableCount)

	must(p.s.ReadObject(&nodes))

	return nodes
}

func (p *Parser) NodeName(n Node) string {
	return p.names[n.NodeNameOffset]
}

func (p *Parser) NodeContent(n Node) string {
	content, ok := p.names[n.ContentOffset]
	if !ok {
		content = "BUGGED"
	}

	return content
}

func (p *Parser) GetAttr(id int32) (name, value string) {
	value, ok := p.names[p.attrs[id].ValueOffset]
	if !ok {
		value = "BUGGED"
	}

	name = p.names[p.attrs[id].NameOffset]

	return name, value
}

func readDict(s *stream.Stream, offset int64) DataMap {
	_, _ = s.Seek(offset, io.SeekStart)

	dataMap := make(DataMap)

	for s.Pos() < s.Size() {
		dataMap[int32(s.Pos()-offset)] = s.ReadCString()
	}

	return dataMap
}

func readAttrs(s *stream.Stream, offset int64, count int) []Reference {
	_, _ = s.Seek(offset, io.SeekStart)

	attrs := make([]Reference, count)

	must(s.ReadObject(&attrs))

	return attrs
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func mustRead[T any](s *stream.Stream) T {
	var v T

	must(s.ReadObject(&v))

	return v
}
