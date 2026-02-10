package replication

import "fmt"

// LogEntry is an entry in the replication log.
type LogEntry struct {
	Key   []byte
	Value []byte
}

// Replicator is the interface for the replication mechanism.
type Replicator interface {
	Replicate(entry *LogEntry) error
}

// Follower is a follower in the replication group.
type Follower struct {
	id string
}

// NewFollower creates a new follower.
func NewFollower(id string) *Follower {
	return &Follower{id: id}
}

// Replicate replicates a log entry.
func (f *Follower) Replicate(entry *LogEntry) error {
	fmt.Printf("Follower %s: Replicating key %s\n", f.id, string(entry.Key))
	// In a real implementation, this would apply the change to the local storage.
	return nil
}

// Leader is the leader in the replication group.
type Leader struct {
	followers []*Follower
}

// NewLeader creates a new leader.
func NewLeader(followers []*Follower) *Leader {
	return &Leader{followers: followers}
}

// Replicate replicates a log entry to all followers.
func (l *Leader) Replicate(entry *LogEntry) error {
	fmt.Printf("Leader: Replicating key %s to %d followers\n", string(entry.Key), len(l.followers))
	for _, follower := range l.followers {
		if err := follower.Replicate(entry); err != nil {
			return err
		}
	}
	return nil
}