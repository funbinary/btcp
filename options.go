package btcp

import (
	"btcp/bpack"
	"btcp/logger"
	"time"
)

type Option func(s *Server)

func WithPacket(packet bpack.IPack) Option {
	return func(s *Server) {
		s.packet = packet
	}
}

func WithSessionManager(need bool) Option {
	return func(s *Server) {
		s.needSessionMgr = need
	}
}

// WithHeartBeatDuration 设置心跳超时时间,单位秒
func WithHeartBeatDuration(d int) Option {
	return func(s *Server) {
		s.heartbeatMaxDuration = time.Duration(d) * time.Second
	}
}

func WithListenIp(ip string) Option {
	return func(s *Server) {
		s.ip = ip
	}
}

func WithPort(p int) Option {
	return func(s *Server) {
		s.port = p
	}
}

func WithTls(crt, key string) Option {
	return func(s *Server) {
		s.certFile = crt
		s.privateKeyFile = key
	}
}

func WithLogger(l logger.ILogger) Option {
	return func(s *Server) {
		s.logger = l
	}
}
