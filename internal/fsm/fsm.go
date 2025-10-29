package fsm

import (
	"sync"
	"time"

	"go_scripts/internal/logger"
)

type StateType string

const (
	StateDefault   StateType = "default"
	StateAwaitLink StateType = "await_link"
)

type State struct {
	Type StateType
	Data map[string]interface{}
}

func NewState(t StateType, data map[string]interface{}) State {
	if data == nil {
		data = map[string]interface{}{}
	}
	return State{Type: t, Data: data}
}

func Start() State     { return NewState(StateDefault, nil) }
func AwaitLink() State { return NewState(StateAwaitLink, nil) }

type UserStateEntry struct {
	State     State
	LastSeen  time.Time
	CreatedAt time.Time
}

type Manager struct {
	states   map[int]*UserStateEntry
	mutex    sync.RWMutex
	ttl      time.Duration
	cleanup  time.Duration
	stopChan chan struct{}
}

func NewManager(ttl, cleanupInterval time.Duration) *Manager {
	m := &Manager{states: map[int]*UserStateEntry{}, ttl: ttl, cleanup: cleanupInterval, stopChan: make(chan struct{})}
	go m.cleanupLoop()
	return m
}

func (m *Manager) Get(userID int) (State, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	e, ok := m.states[userID]
	if !ok {
		return State{}, false
	}
	e.LastSeen = time.Now()
	return e.State, true
}

func (m *Manager) Set(userID int, s State) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	now := time.Now()
	if e, ok := m.states[userID]; ok {
		e.State = s
		e.LastSeen = now
	} else {
		m.states[userID] = &UserStateEntry{State: s, LastSeen: now, CreatedAt: now}
		logger.UserInfo(userID, "Новый пользователь")
	}
}

func (m *Manager) cleanupLoop() {
	t := time.NewTicker(m.cleanup)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			m.cleanupExpiredStates()
		case <-m.stopChan:
			return
		}
	}
}

func (m *Manager) cleanupExpiredStates() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	now := time.Now()
	expired := 0
	for uid, e := range m.states {
		if now.Sub(e.LastSeen) > m.ttl {
			delete(m.states, uid)
			expired++
		}
	}
	if expired > 0 {
		logger.BotInfo("Очищено %d устаревших состояний пользователей", expired)
	}
}

func (m *Manager) Shutdown() { close(m.stopChan); logger.BotInfo("State manager shutdown completed") }
