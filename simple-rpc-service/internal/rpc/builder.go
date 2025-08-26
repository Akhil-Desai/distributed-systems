package stub

import (
	"encoding/binary"
)

type MessageBuilder struct {
	signature string
	a         int32
	b         int32
}

func (mb *MessageBuilder) SetSignature(sig string) *MessageBuilder{
	mb.signature = sig
	return mb
}

func (mb *MessageBuilder) SetA(a int32) *MessageBuilder {
	mb.a = a
	return mb
}

func (mb *MessageBuilder) SetB(b int32) *MessageBuilder {
	mb.b = b
	return mb
}


func (mb *MessageBuilder) Build() ([]byte,error) {

	msg := make([]byte, 1024)
	offset := 0

	binary.BigEndian.PutUint32(msg, uint32(len(mb.signature)))
	offset += 4

	copy(msg[offset:], []byte(mb.signature))
	offset += len(mb.signature)

	binary.BigEndian.PutUint32(msg[offset:], uint32(mb.a))
	offset += 4

	binary.BigEndian.PutUint32(msg[offset:], uint32(mb.b))
	offset += 4

	return msg,nil
}
