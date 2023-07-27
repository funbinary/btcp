package bpack

type Message struct {
	Version uint32
	BizCode uint32
	SeqID   uint64
	DataLen uint32
	Data    []byte
}

func NewMsgPackage(version uint32, code uint32, seqid uint64, data []byte) *Message {
	return &Message{
		Version: version,
		BizCode: code,
		SeqID:   seqid,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

func NewMessage(len uint32, data []byte) *Message {
	return &Message{
		DataLen: len,
		Data:    data,
	}
}

func (msg *Message) Clone() *Message {
	return &Message{
		Version: msg.Version,
		BizCode: msg.BizCode,
		SeqID:   msg.SeqID,
		DataLen: msg.DataLen,
		Data:    msg.Data,
	}
}

func (msg *Message) GetDataLen() uint32 {
	return msg.DataLen
}

func (msg *Message) GetVersion() uint32 {
	return msg.Version
}

func (msg *Message) GetSeqID() uint64 {
	return msg.SeqID
}

func (msg *Message) GetBizCode() uint32 {
	return msg.DataLen
}

func (msg *Message) GetData() []byte {
	return msg.Data
}

func (msg *Message) SetVersion(v uint32) {
	msg.Version = v
}

func (msg *Message) SetBizCode(code uint32) {
	msg.BizCode = code
}

func (msg *Message) SetSeqID(seqId uint64) {
	msg.SeqID = seqId
}

func (msg *Message) SetData(data []byte) {
	msg.Data = data
}

func (msg *Message) SetDataLen(len uint32) {
	msg.DataLen = len
}
