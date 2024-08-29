package node_test

import (
	"context"
	"errors"
	"testing"
	"time"
	"warehouse/internal/raft-cluster/node"
	"warehouse/pkg/raft/raft_cluster_v1"
)

func TemplateNetwork() []*node.ClusterNodeServer {

	network0 := []*node.ClusterNodeServer{}
	network1 := []*node.ClusterNodeServer{}
	network2 := []*node.ClusterNodeServer{}

	leader := node.New(0, node.Lead, 1, network0)
	follower1 := node.New(1, node.Follower, leader.IdNode, network1)
	follower2 := node.New(2, node.Follower, leader.IdNode, network2)

	network0 = append(network0, follower1)
	network0 = append(network0, follower2)

	network1 = append(network1, leader)
	network1 = append(network1, follower2)

	network2 = append(network2, leader)
	network2 = append(network2, follower1)

	leader.Network = network0
	follower1.Network = network1
	follower2.Network = network2

	return []*node.ClusterNodeServer{leader, follower1, follower2}
}

func TestSimpleHeartBeat(t *testing.T) {
	cluster := TemplateNetwork()
	lead := cluster[0]

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*5))
	defer cancel()

	_, err := lead.SendHeartBeat(ctx, &raft_cluster_v1.Empty{})

	if !errors.Is(err, context.Canceled) {
		t.Errorf(err.Error())
	}
}
