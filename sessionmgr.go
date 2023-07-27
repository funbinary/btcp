package btcp

import (
	"sync"
)

type SessionManager struct {
	sessions sync.Map
}

func newSessionManager() *SessionManager {
	return &SessionManager{}
}

func (sm *SessionManager) Add(sessID string, sess ISession) {
	sm.sessions.Store(sessID, sess)
}

func (sm *SessionManager) Remove(sessID string) {
	sm.sessions.Delete(sessID)
}

func (sm *SessionManager) Find(sessID string) ISession {
	session, ok := sm.sessions.Load(sessID)
	if ok {
		return session.(ISession)
	}
	return nil
}

func (sm *SessionManager) Range(f func(sessionId, session any) bool) {
	sm.sessions.Range(func(key, value any) bool {
		return f(key, value)
	})
}
