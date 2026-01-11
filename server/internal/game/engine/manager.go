package engine

import (
    "sync"
)

// Global instance (Thread-safe map)
var GlobalManager = &Manager{
    sessions: make(map[string]*Session),
}

type Manager struct {
    sessions map[string]*Session
    lock     sync.RWMutex
}

func (m *Manager) AddSession(s *Session) {
    m.lock.Lock()
    defer m.lock.Unlock()
    m.sessions[s.UserID] = s
}

func (m *Manager) GetSession(userID string) *Session {
    m.lock.RLock()
    defer m.lock.RUnlock()
    return m.sessions[userID]
}
