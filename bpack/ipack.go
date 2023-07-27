package bpack

type IPack interface {
	GetHeadLen() uint32                // 获取包头长度方法
	Pack(msg IMessage) ([]byte, error) // 封包方法
	Unpack([]byte) []IMessage          // 拆包方法
}
