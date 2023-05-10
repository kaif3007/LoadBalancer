package manager

import (
	"fmt"
	"loadbalancer/APIRequest"
	"sync"
)

// APINameToServerIPMapping is a mapping for APIName to serverIPs serving it
// lock is there to protect it from concurrent access
type APINameToServerIPMapping struct {
	mapping map[string][]string
	sync.Mutex
}

type Manager struct {
	serverIPMapping *APINameToServerIPMapping
	matcher         MatcherIntf
	rateLimiter     RateLimiterIntf
}

func NewManager(matcherType string, rateLimiterType string) *Manager {
	m := &Manager{
		serverIPMapping: &APINameToServerIPMapping{
			mapping: make(map[string][]string),
		},
	}

	m.ConfigureMatcher(matcherType)
	m.ConfigureRateLimiter(rateLimiterType)
	return m
}

func (m *Manager) ConfigureMatcher(matcherType string) {
	switch matcherType {
	case "ROUND_ROBIN":
		// other types of matchers can also be configured
		m.matcher = NewRoundRobin(m)
	}
}

func (m *Manager) ConfigureRateLimiter(rateLimiterType string) {
	switch rateLimiterType {
	case "TOKEN_BUCKET":
		// other types of rate limiters can also be configured
		m.rateLimiter = NewTokenBucket(10, 1)
	}
}

// GetServerIP returns IP of server given APIName and index
func (m *Manager) GetServerIP(APIName string, idx int) (string, error) {
	m.serverIPMapping.Lock()
	defer m.serverIPMapping.Unlock()

	serverIPs := m.serverIPMapping.mapping[APIName]
	if idx >= len(serverIPs) {
		return "", fmt.Errorf("server IP at index %v, does not exists", idx)
	}

	return serverIPs[idx], nil
}

// GetNumServers returns no. of servers available for given API
func (m *Manager) GetNumServers(APIName string) int {
	m.serverIPMapping.Lock()
	defer m.serverIPMapping.Unlock()

	serverIPs := m.serverIPMapping.mapping[APIName]
	return len(serverIPs)
}

// ApplyConfigs is used to dynamically change serverIPs for given API Name
func (m *Manager) ApplyConfigs(APIName string, serverIPs []string) {
	m.serverIPMapping.Lock()
	defer m.serverIPMapping.Unlock()

	m.serverIPMapping.mapping[APIName] = serverIPs
	m.matcher.Reset()
}

// HandleRequest is called for each request, first its checked by rate limiter
// and then its destination IP is changed to serverIP matched by Loadbalancer
func (m *Manager) HandleRequest(req *APIRequest.APIRequest) error {
	if !m.rateLimiter.IsRequestAllowed(req) {
		return fmt.Errorf("not allowed to make too many requests")
	}

	return m.matcher.AssignServerIP(req)
}
