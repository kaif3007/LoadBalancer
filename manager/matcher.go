package manager

import (
	"fmt"
	"loadbalancer/APIRequest"
	"sync"
)

// MatcherIntf is supposed to be implemented by any matching algorithm
type MatcherIntf interface {
	// AssignServerIP - given an API request, change its destinationIP to one
	// of the available servers serving the API
	AssignServerIP(req *APIRequest.APIRequest) error
	// Reset() is supposed to be called by manager in case of any config changes
	// In the method matcher should clear all state its managing
	Reset()
}

// RoundRobin implements MatcherIntf interface
type RoundRobin struct {
	manager            *Manager
	APINameServerIndex map[string]int
	lock               sync.Mutex
}

func NewRoundRobin(m *Manager) *RoundRobin {
	return &RoundRobin{
		manager:            m,
		APINameServerIndex: make(map[string]int),
	}
}

func (r *RoundRobin) Reset() {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.APINameServerIndex = make(map[string]int)
}

func (r *RoundRobin) AssignServerIP(req *APIRequest.APIRequest) error {
	APIName := req.Name
	numServers := r.manager.GetNumServers(APIName)
	if numServers == 0 {
		return fmt.Errorf("no servers available for APIName %v", APIName)
	}

	r.lock.Lock()
	currentServerIdx := r.APINameServerIndex[APIName]
	r.APINameServerIndex[APIName]++
	r.APINameServerIndex[APIName] %= numServers
	r.lock.Unlock()

	serverIP, err := r.manager.GetServerIP(APIName, int(currentServerIdx))
	if err != nil {
		return err
	}

	req.DestinationIP = serverIP
	return nil
}
