package events

import (
	"sync"
)

func AddMessage(sessionID string, msg *Message) {
	mutex.Lock()

	m, ok := messages[sessionID]
	if !ok {
		m = make([]*Message, 0)
	}
	m = append(m, msg)
	messages[sessionID] = m

	c, ok2 := sessions[sessionID]

	mutex.Unlock()

	if ok2 && c != nil {
		c <- struct{}{}
	}
}

func GetMessages(sessionID string) []*Message {
	mutex.RLock()
	defer mutex.RUnlock()
	return messages[sessionID]
}

func CommitMessages(sessionID string, msgs []*Message) {
	l1 := len(msgs)

	mutex.Lock()
	if m, ok := messages[sessionID]; ok {
		if l2 := len(m); l2 <= l1 {
			m = m[:0]
		} else {
			m = m[:l1-l2]
		}
		messages[sessionID] = m
	}
	mutex.Unlock()
}

func AddSession(sessionID string, trigger chan<- struct{}) {
	mutex.Lock()
	sessions[sessionID] = trigger
	if _, ok := messages[sessionID]; !ok {
		messages[sessionID] = make([]*Message, 0)
	}
	mutex.Unlock()
}

func RemoveSession(sessionID string) {
	mutex.Lock()
	delete(sessions, sessionID)
	mutex.Unlock()
}

var (
	mutex    sync.RWMutex
	sessions = make(map[string]chan<- struct{}, 128)
	messages = make(map[string][]*Message, 128)
)
