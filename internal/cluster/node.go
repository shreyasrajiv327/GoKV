package cluster

import (
	"sync"
    "time"
)

type Role string

const(
	Leader Role = "leader"
	Follower Role = "follower"
	Candidate Role = "candidate"
)

type Node struct{

	mu sync.RWMutex

	ID string
	Role Role

	CurrentTerm int

	LastHeartbeat time.Time
}

func NewNode(id string, role Role) *Node{
	return &Node{
		ID: id,

		Role: role,

		CurrentTerm: 0,

		LastHeartbeat: time.Now(),
	}
}


func (n *Node) SetRole(role Role){
	n.mu.Lock()
	defer n.mu.Unlock()

	n.Role = role
}

func (n *Node) GetRole() Role {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return n.Role
}

func (n *Node) IncrementTerm(){
	n.mu.Lock()
	defer n.mu.Unlock()

	n.CurrentTerm++
}

func (n *Node) GetTerm() int{
	n.mu.RLock()
	defer n.mu.RUnlock()

	return n.CurrentTerm
}

func (n *Node) UpdateHeartbeat(){
	n.mu.Lock()
	defer n.mu.Unlock()

	n.LastHeartbeat = time.Now()
}

func(n *Node) TimeSinceHeartbeat() time.Duration{
	n.mu.RLock()
	defer n.mu.RUnlock()

	return time.Since(n.LastHeartbeat)
}

