package node

import (
	"fmt"
	"warehouse/internal/model"
)

const (
	Follower = iota
	Candidate
	Lead
)

type ClusterNode struct {
	idNode int	
	status   bool
	settings *config

	log []model.Instance

	idLead int
}

func NewClusterNode(id int) *ClusterNode {

	return &ClusterNode{
		idNode: id,
		status: true,
		settings: load(),

		log: make([]model.Instance, 0),

	}
}

func (r *ClusterNode) Serve() {}

// send HeartBeat string to global monitoring system
func (r *ClusterNode) GetHeartBeat() string {
	res := fmt.Sprintf("%d", r.settings.Replication) + r.settings.Port
	if !r.status {
		res += " stopped"
	}
	return res
}
