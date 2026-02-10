package cluster

import (
	"github.com/lafikl/consistent"
)

// Node represents a node in the cluster.
type Node struct {
	ID   string
	Addr string
}

// Cluster is the interface for the cluster manager.
type Cluster interface {
	Nodes() ([]*Node, error)
	AddNode(node *Node) error
	RemoveNode(node *Node) error
	GetNode(key string) (*Node, error)
}

// Manager is the implementation of the cluster manager.
type Manager struct {
	nodes []*Node
	ring  *consistent.Consistent
}

// NewManager creates a new Manager.
func NewManager(nodes []*Node) *Manager {
	c := consistent.New()
	for _, node := range nodes {
		c.Add(node.ID)
	}
	return &Manager{
		nodes: nodes,
		ring:  c,
	}
}

// Nodes returns the nodes in the cluster.
func (m *Manager) Nodes() ([]*Node, error) {
	return m.nodes, nil
}

// AddNode adds a node to the cluster.
func (m *Manager) AddNode(node *Node) error {
	m.nodes = append(m.nodes, node)
	m.ring.Add(node.ID)
	return nil
}

// RemoveNode removes a node from the cluster.
func (m *Manager) RemoveNode(node *Node) error {
	for i, n := range m.nodes {
		if n.ID == node.ID {
			m.nodes = append(m.nodes[:i], m.nodes[i+1:]...)
			break
		}
	}
	m.ring.Remove(node.ID)
	return nil
}

// GetNode gets the node for a given key.
func (m *Manager) GetNode(key string) (*Node, error) {
	nodeID, err := m.ring.Get(key)
	if err != nil {
		return nil, err
	}
	for _, n := range m.nodes {
		if n.ID == nodeID {
			return n, nil
		}
	}
	return nil, nil
}
