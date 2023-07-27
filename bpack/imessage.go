package bpack

type IMessage interface {
	GetDataLen() uint32 // 获取消息数据段长度
	GetVersion() uint32 // 获取版本号
	GetSeqID() uint64   // 获取消息ID
	GetBizCode() uint32 // 获取业务码
	GetData() []byte    // 获取消息内容

	SetVersion(v uint32)    // 设置版本号
	SetBizCode(code uint32) // 设置业务码
	SetSeqID(uint64)        // 设置消息ID
	SetData([]byte)         // 设置消息内容
	SetDataLen(uint32)      // 设置消息数据段长度
}
