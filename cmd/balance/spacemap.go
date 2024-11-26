package main

import (
	"sync"
)

// All space values are in int64(bytes)
type DBMap struct {
	m *sync.RWMutex
	// This maps from the hostname to the server in the
	// infra. E.g. map["www.va.gov"]"s0"
	HostToServer map[string]string
	// This tells us how much space has been used
	// by a given server. E.g. map["s0"]123456
	ServerUsed map[string]int64
	// In case it ever varies...
	// E.g. map["s0"]5000000000
	ServerMax map[string]int64
}

func NewDBMap(max int64) *DBMap {
	return &DBMap{
		m:            &sync.RWMutex{},
		HostToServer: make(map[string]string),
		ServerUsed:   make(map[string]int64),
		ServerMax:    make(map[string]int64),
	}
}

// The servers we know about. Should be dynamic/based on
// config instead of coded here.
var Servers []string = []string{"s0", "s1"}
// This is an index; we bump by one and % length, so that
// we start with a new server with each addition.
var NextStart int = 0

// For now, adding a server is round-robin. The right way would be
// to consider reshuffling all of the servers with every addition, and 
// issue change messages to the queues. This is less space efficient/not
// guaranteed to be optimal.
func (sm *DBMap) AddDB(hostname string) error {
	searching := true
	for searching {
		if 
	}
}

// This overwrites anything we have. It updates the space information
// for a server. They do this periodically.
func (sm *DBMap) UpsertServer(server string, serverUsed, maxAvail int64) {
	sm.m.Lock()
	defer sm.m.Unlock()
	sm.ServerUsed[server] = serverUsed
	sm.ServerMax[server] = maxAvail
}

func (sm *DBMap) Remove(server string) {
	sm.m.Lock()
	defer sm.m.Unlock()
	delete(sm.ServerUsed, server)
	delete(sm.ServerMax, server)
}
