package routing

import "github.com/gemini-cli/distributed-storage-engine/internal/cluster"

// QueryRouter is the interface for the query router.
type QueryRouter interface {
	Route(key string) (*cluster.Node, error)
}

// Router is the implementation of the query router.
type Router struct {
	cluster cluster.Cluster
}

// NewRouter creates a new Router.
func NewRouter(cluster cluster.Cluster) *Router {
	return &Router{cluster: cluster}
}

// Route routes a request to the correct node.
func (r *Router) Route(key string) (*cluster.Node, error) {
	return r.cluster.GetNode(key)
}
