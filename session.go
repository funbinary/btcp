package btcp

import (
	"btcp/bpack"
	"btcp/logger"
	"context"
	"encoding/hex"
	"errors"
	"net"
	"strconv"
	"sync"
	"time"
)

type ISession interface {
	Start()
	SessionId() uint64
}

type Session struct {
	conn                 net.Conn
	heartbeatMaxDuration time.Duration
	ctx                  context.Context
	cancel               context.CancelFunc
	root                 *RouterGroup

	logger logger.ILogger // 日志,默认使用fmt

	connID uint64

	lastActivityTime time.Time // 最后一次活动时间
	isClosed         bool      // 当前连接的关闭状态

	property map[string]interface{} // 属性

	onConnStart func(sess ISession) // Server的连接创建时Hook函数
	onConnStop  func(sess ISession) // Server的连接断开时的Hook函数
	sessionMgr  *SessionManager
	packet      bpack.IPack
	rwLock      sync.RWMutex
}

func (sess *Session) Start() {
	// 启用心跳检测
	defer func() {
		if err := recover(); err != nil {
			sess.logger.Errorf("Connection Start() error: %v", err)
		}
	}()

	// Execute the hook method for processing business logic when creating a connection
	// (按照用户传递进来的创建连接时需要处理的业务，执行钩子方法)
	//sess.callOnConnStart()

	// Start heartbeating detection
	//if sess.hc != nil {
	//	sess.hc.Start()
	//	sess.updateActivity()
	//}

	// Start the Goroutine for reading data from the client
	// (开启用户从客户端读取数据流程的Goroutine)
	go sess.reader()

	select {
	case <-sess.ctx.Done():
		sess.finalizer()

		// 归还workerid
		//freeWorker(sess)
		return
	}
}

func (sess *Session) Stop() {
	sess.cancel()
}

func (sess *Session) Send(data []byte) error {
	sess.rwLock.RLock()
	defer sess.rwLock.RUnlock()
	if sess.isClosed == true {
		return errors.New("connection closed when send msg")
	}

	_, err := sess.conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (sess *Session) IsAlive() bool {
	if sess.isClosed {
		return false
	}
	// 检查连接最后一次活动时间，如果超过心跳间隔，则认为连接已经死亡
	return time.Now().Sub(sess.lastActivityTime) < sess.heartbeatMaxDuration
}

func (sess *Session) updateActivity() {
	sess.lastActivityTime = time.Now()
}

func (sess *Session) SessionId() uint64 {
	return sess.connID
}

func (sess *Session) SessionIdStr() string {
	return strconv.FormatUint(sess.connID, 10)
}

func (sess *Session) RemoteAddr() net.Addr {
	return sess.conn.RemoteAddr()
}

func (sess *Session) LocalAddr() net.Addr {
	return sess.conn.LocalAddr()
}

func (sess *Session) reader() {
	sess.logger.Infof("%s [Reader Goroutine is running]", sess.RemoteAddr().String())
	defer sess.logger.Infof("%s [conn Reader exit!]", sess.RemoteAddr().String())
	defer sess.Stop()
	defer func() {
		if err := recover(); err != nil {
			sess.logger.Errorf("connID=%d, panic err=%v", sess.SessionId(), err)
		}
	}()
	buffer := make([]byte, 1024)
	for {
		select {
		case <-sess.ctx.Done():
			return
		default:
			// (从conn的IO中读取数据到内存缓冲buffer中)
			n, err := sess.conn.Read(buffer)
			if err != nil {
				sess.logger.Errorf("read msg head [read datalen=%d], error = %s", n, err)
				return
			}
			sess.logger.Tracef("read buffer %s \n", hex.EncodeToString(buffer[0:n]))

			// 更新心跳时间
			sess.updateActivity()

			if sess.packet != nil {
				msgs := sess.packet.Unpack(buffer[:n])
				if len(msgs) <= 0 {
					continue
				}

				for _, msg := range msgs {
					req := sess.newRequest(msg)
					req.Next()
				}
			} else {
				sess.logger.Fatalf("No Found packet")
			}
		}
	}
}

func (sess *Session) newRequest(msg bpack.IMessage) *Request {
	handlers, ok :=sess.root.Apis[msg.GetBizCode()]
	if !ok {

	}
	return &Request{
		session: sess,
		msg:     msg,
		index:   -1,
		handlers:
	}
}

func (sess *Session) finalizer() {
	// Call the callback function registered by the user when closing the connection if it exists
	// (如果用户注册了该链接的	关闭回调业务，那么在此刻应该显示调用)
	sess.onConnStop(sess)

	sess.rwLock.Lock()
	defer sess.rwLock.Unlock()

	// If the connection has already been closed
	if sess.isClosed == true {
		return
	}

	// Stop the heartbeat detector associated with the connection
	//if sess.hc != nil {
	//	sess.hc.Stop()
	//}
	_ = sess.conn.Close()

	// Remove the connection from the connection manager
	//if sess.connManager != nil {
	//	sess.connManager.Remove(sess)
	//}

	// Close all channels associated with the connection
	//if sess.msgBuffChan != nil {
	//	close(sess.msgBuffChan)
	//}

	sess.isClosed = true

	//zlog.Ins().InfoF("Conn Stop()...ConnID = %d", sess.connID)
}
