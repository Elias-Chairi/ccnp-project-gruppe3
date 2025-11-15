package tcpservice

import (
	"sync"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util"
)

type PendingRequests struct {
	mu   sync.RWMutex
	reqs map[uint16]messages.TopLevelMessage
}

func NewPendingRequests() *PendingRequests {
	return &PendingRequests{
		reqs: make(map[uint16]messages.TopLevelMessage),
	}
}

func (p *PendingRequests) Add(msg messages.TopLevelMessage) uint16 {
	p.mu.Lock()
	defer p.mu.Unlock()
	key := util.GetUniqueMapKey(p.reqs)
	p.reqs[key] = msg
	return key
}

func (p *PendingRequests) Get(reqID uint16) (messages.TopLevelMessage, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	msg, exists := p.reqs[reqID]
	return msg, exists
}

func (p *PendingRequests) Remove(reqID uint16) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.reqs, reqID)
}
