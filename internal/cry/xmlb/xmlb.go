package xmlb

import (
	"github.com/beevik/etree"

	"github.com/GlebYaltchik/sc-keybind-extract/internal/stream"
)

// This code is a porting of the CryEngine XMLB parser to Go.
// The original code is written in C# and can be found here:
// https://github.com/dolkensp/unp4k/blob/develop/src/unforge/CryXmlB/CryXmlSerializer.cs

func Decode(data []byte) ([]byte, error) {
	if len(data) > 0 && data[0] == '<' {
		return data, nil // not encoded
	}

	br := stream.New(data)

	p, err := NewParser(br)
	if err != nil {
		return nil, err
	}

	doc := etree.NewDocument()

	attributeIndex := int32(0)
	xmlMap := make(map[int]*etree.Element)

	for nodeID, node := range p.ReadNodes() {
		element := etree.NewElement(p.NodeName(node))

		for i := int16(0); i < node.AttributeCount; i++ {
			element.CreateAttr(p.GetAttr(attributeIndex))
			attributeIndex++
		}

		xmlMap[nodeID] = element

		if content := p.NodeContent(node); content != "" {
			element.AddChild(etree.NewCData(content))
		}

		parent, ok := xmlMap[int(node.ParentNodeID)]
		if ok {
			parent.AddChild(element)
		} else {
			doc.AddChild(element)
		}
	}

	doc.Indent(2)

	return doc.WriteToBytes()
}
