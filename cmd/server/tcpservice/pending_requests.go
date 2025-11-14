package tcpservice

import (
	"sync"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
)

type PendingRequests struct {
	mu    sync.RWMutex
	nodes map[uint8]encoding.TLV
}

func NewPendingRequests() *PendingRequests {
	return &PendingRequests{
		nodes: make(map[uint8]encoding.TLV),
	}
}

func (pr *PendingRequests) Add(nodeID uint8, tlv encoding.TLV) {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	pr.nodes[nodeID] = tlv
}

func (pr *PendingRequests) Get(nodeID uint8) (encoding.TLV, bool) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()
	tlv, exists := pr.nodes[nodeID]
	return tlv, exists
}

func (pr *PendingRequests) Remove(nodeID uint8) {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	delete(pr.nodes, nodeID)
}
