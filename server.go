package btcp

import (
	"btcp/bpack"
	"btcp/logger"
	"context"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

var (
	defaultPort                 int = 20000
	defaultHeartbeatMaxDuration     = time.Duration(10) * time.Second
)

type IServer interface {
}

type Server struct {
	root *RouterGroup
	ip   string // 监听的IP地址,空为监听全部
	port int    // 监听的端口,默认为20000

	certFile       string // tls certificate file.
	privateKeyFile string // tls private key file.

	cID                  uint64
	packet               bpack.IPack
	heartbeatMaxDuration time.Duration
	ctx                  context.Context
	cancel               context.CancelFunc
	logger               logger.ILogger      // 日志,默认使用fmt
	onConnStart          func(sess ISession) // 该Server的连接创建时Hook函数
	onConnStop           func(sess ISession) // Server的连接断开时的Hook函数
	sessionMgr           *SessionManager
	needSessionMgr       bool //是否需要session管理,否为业务层手动管理
}

func NewServer(opts ...Option) IServer {
	s := &Server{
		heartbeatMaxDuration: defaultHeartbeatMaxDuration,
		logger:               &logger.Logger{},
		ip:                   "",
		port:                 defaultPort,
		sessionMgr:           newSessionManager(),
		packet:               bpack.NewPack(),
		needSessionMgr:       true,
		root:                 NewRouteGroup(),
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())

	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) Run() error {
	network := "tcp"
	address := fmt.Sprintf("%s:%d", s.ip, s.port)
	addr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return err
	}

	// 2. Listen to the server address
	var listener net.Listener
	if s.certFile != "" && s.privateKeyFile != "" {
		// Read certificate and private key
		crt, err := tls.LoadX509KeyPair(s.certFile, s.privateKeyFile)
		if err != nil {
			return err
		}

		// TLS connection
		tlsConfig := &tls.Config{}
		tlsConfig.Certificates = []tls.Certificate{crt}
		tlsConfig.Time = time.Now
		tlsConfig.Rand = rand.Reader
		listener, err = tls.Listen(network, address, tlsConfig)
		if err != nil {
			return err
		}
	} else {
		listener, err = net.ListenTCP(network, addr)
		if err != nil {
			return err
		}
	}

	for {
		// 阻塞等待客户端建立连接请求
		conn, err := listener.Accept()
		if err != nil {
			//Go 1.16+
			if errors.Is(err, net.ErrClosed) {
				s.logger.Errorf("Listener closed")
				return err
			}
			s.logger.Errorf("Accept err: %v", err)
			continue
		}

		// 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
		newCid := atomic.AddUint64(&s.cID, 1)
		sess := s.newSession(conn, newCid)

		go sess.Start()

	}
	return nil
}

func (s *Server) newSession(conn net.Conn, connID uint64) ISession {
	// Initialize Conn properties
	sess := &Session{
		conn:                 conn,
		heartbeatMaxDuration: s.heartbeatMaxDuration,
		ctx:                  nil,
		cancel:               nil,
		logger:               s.logger,
		connID:               connID,
		isClosed:             false,
		property:             nil,
	}

	sess.ctx, sess.cancel = context.WithCancel(s.ctx)

	sess.onConnStart = s.onConnStart
	sess.onConnStop = s.onConnStop
	sess.sessionMgr = s.sessionMgr
	sess.packet = s.packet
	sess.root = s.root

	if s.needSessionMgr {
		sess.sessionMgr.Add(sess.SessionIdStr(), sess)
	}

	return sess
}

func (s *Server) SetLogger(l logger.ILogger) {
	s.logger = l
}

func (s *Server) SetOnConnStart(hookFunc func(session ISession)) {
	s.onConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(session ISession)) {
	s.onConnStop = hookFunc
}

func (s *Server) GetOnConnStart() func(ISession) {
	return s.onConnStart
}

func (s *Server) GetOnConnStop() func(ISession) {
	return s.onConnStop
}

func (s *Server) Use(handlers ...HandlerFunc) IRoutes {
	return s.root.Use(handlers...)
}

func (s *Server) Handle(code uint32, handlers ...HandlerFunc) IRoutes {
	return s.root.Handle(code, handlers...)
}
