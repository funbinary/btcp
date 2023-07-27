package bpack

import (
	"bytes"
	"encoding/binary"
	"sync"
)

var defaultHeaderLen uint32 = 20

type BPack struct {
	buffer *bytes.Buffer
	cache  *Message
	lock   sync.Mutex
}

func NewPack() IPack {
	return &BPack{
		buffer: bytes.NewBuffer([]byte{}),
	}
}

func (p *BPack) GetHeadLen() uint32 {
	// 版本号+包体长度+业务码+序列ID
	return defaultHeaderLen
}

func (p *BPack) Pack(msg IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 版本号
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetVersion()); err != nil {
		return nil, err
	}

	// 包体长度
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// 业务码
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetBizCode()); err != nil {
		return nil, err
	}

	// 序列ID
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetSeqID()); err != nil {
		return nil, err
	}

	// 数据
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (p *BPack) Unpack(data []byte) []IMessage {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.buffer.Write(data)
	var msgs []IMessage

	for {
		if p.cache == nil {
			if p.buffer.Len() < int(defaultHeaderLen) {
				// 长度不够
				return msgs
			}

			msg := &Message{}

			// 版本号
			if err := binary.Read(p.buffer, binary.BigEndian, &msg.Version); err != nil {
				return msgs
			}

			// 包体长度
			if err := binary.Read(p.buffer, binary.BigEndian, &msg.DataLen); err != nil {
				return msgs
			}

			// 业务码
			if err := binary.Read(p.buffer, binary.BigEndian, &msg.BizCode); err != nil {
				return msgs
			}

			// 序列ID
			if err := binary.Read(p.buffer, binary.BigEndian, &msg.SeqID); err != nil {
				return msgs
			}
			if uint32(p.buffer.Len()) < msg.GetDataLen() {
				p.cache = msg
				return msgs
			}
			msg.SetData(p.buffer.Next(int(msg.GetDataLen())))
			msgs = append(msgs, msg)
		} else {
			if uint32(p.buffer.Len()) < p.cache.GetDataLen() {
				// 长度不够
				return msgs
			}
			msg := p.cache.Clone()
			msg.SetData(p.buffer.Next(int(msg.GetDataLen())))
			msgs = append(msgs, msg)
			p.cache = nil
		}
	}
	return msgs
}
