package node_test

import (
	"context"
	"errors"
	"testing"
	"warehouse/internal/model"
	"warehouse/internal/raft-cluster/node"
	"warehouse/pkg/raft/raft_cluster_v1"

	"github.com/stretchr/testify/assert"
)

func newTestNode(id int, state node.StateType, term int64, logs []model.Instance) *node.ClusterNodeServer {
	return &node.ClusterNodeServer{
		IdNode:   id,
		Term:     term,
		State:    state,
		SizeLogs: len(logs),
		Logs:     logs,
	}
}

func TestStartElection_StateNotFollower(t *testing.T) {
	testNode := newTestNode(1, node.Lead, 100, nil)

	resp, err := testNode.StartElection(context.Background(), &raft_cluster_v1.Empty{})
	assert.Equal(t, "state unsuppoted to election", err.Error())
	assert.Equal(t, int64(100), resp.Term)
}

func TestStartElection_ContextDone(t *testing.T) {
	testNode := newTestNode(1, node.Follower, 100, nil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	resp, err := testNode.StartElection(ctx, &raft_cluster_v1.Empty{})
	assert.True(t, errors.Is(err, context.Canceled))
	assert.Equal(t, int64(100), resp.Term)
}

func TestStartElection_QuorumSatisfied(t *testing.T) {
	// Создаем сеть узлов, где все узлы поддерживают голосование
	node1 := newTestNode(0, node.Follower, 100, []model.Instance{newLogInstance(100, "log1")})
	node2 := newTestNode(1, node.Follower, 100, []model.Instance{newLogInstance(100, "log1")})
	node3 := newTestNode(2, node.Follower, 100, []model.Instance{newLogInstance(100, "log1")})
	node4 := newTestNode(3, node.Follower, 100, []model.Instance{newLogInstance(100, "log1")})
	node5 := newTestNode(4, node.Follower, 100, []model.Instance{newLogInstance(100, "log1")})

	// Сеть
	node1.Network = []*node.ClusterNodeServer{node2, node3, node4, node5}
	node2.Network = []*node.ClusterNodeServer{node1, node3, node4, node5}
	node3.Network = []*node.ClusterNodeServer{node1, node2, node4, node5}
	node4.Network = []*node.ClusterNodeServer{node1, node2, node3, node5}
	node5.Network = []*node.ClusterNodeServer{node1, node2, node3, node4}
	// Старт выборов на узле 1
	resp, err := node1.StartElection(context.Background(), &raft_cluster_v1.Empty{})
	assert.NoError(t, err)
	assert.Equal(t, int64(101), resp.Term)
	assert.Equal(t, node.Lead, node1.State)
	assert.Equal(t, -1, node1.LeadId)
}

func TestStartElection_DisconnectElection(t *testing.T) {
	// Создаем сеть узлов, где один узел не поддерживает голосование
	node1 := newTestNode(0, node.Follower, 100, []model.Instance{newLogInstance(100, "log1")})
	node2 := newTestNode(1, node.Follower, 101, []model.Instance{newLogInstance(100, "log1")})
	node3 := newTestNode(2, node.Follower, 101, []model.Instance{newLogInstance(100, "log1")})
	node4 := newTestNode(3, node.Follower, 100, []model.Instance{newLogInstance(100, "log1")})
	node5 := newTestNode(4, node.Follower, 100, []model.Instance{newLogInstance(100, "log1")})

	// Узел 3 отказывается голосовать
	node3.State = node.Lead // Узел 3 уже лидер и не будет голосовать за выборы

	// Сеть
	node1.Network = []*node.ClusterNodeServer{node2, node3, node4, node5}
	node2.Network = []*node.ClusterNodeServer{node1, node3, node4, node5}
	node3.Network = []*node.ClusterNodeServer{node1, node2, node4, node5}
	node4.Network = []*node.ClusterNodeServer{node1, node2, node3, node5}
	node5.Network = []*node.ClusterNodeServer{node1, node2, node3, node4}

	// Старт выборов на узле 1
	resp, err := node1.StartElection(context.Background(), &raft_cluster_v1.Empty{})
	assert.Equal(t, "voter's logs more complete, candidate not legitimate", err.Error())
	assert.Equal(t, int64(101), resp.Term)
	assert.Equal(t, node.Follower, node1.State)
}

func TestStartElection_QuorumNotSatisfied(t *testing.T) {
	node1 := newTestNode(0, node.Follower, 100, []model.Instance{newLogInstance(100, "log1")})

	node1.Network = []*node.ClusterNodeServer{}

	resp, err := node1.StartElection(context.Background(), &raft_cluster_v1.Empty{})
	assert.Equal(t, "qourum not satisfied", err.Error())
	assert.Equal(t, int64(101), resp.Term)
	assert.Equal(t, node.Follower, node1.State)
}