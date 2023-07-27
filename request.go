package btcp

import "btcp/bpack"

type Request struct {
	session  ISession
	msg      bpack.IMessage
	handlers HandlersChain
	index    int8
}

func (r *Request) GetSession() ISession {
	return r.session
}

func (r *Request) GetSeqID() uint64 {
	return r.msg.GetSeqID()
}

func (r *Request) GetBizCode() uint32 {
	return r.msg.GetBizCode()
}

func (r *Request) GetVersion() uint32 {
	return r.msg.GetVersion()
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) Next() {
	r.index++
	for r.index < int8(len(r.handlers)) {
		r.handlers[r.index](r)
		r.index++
	}
}
