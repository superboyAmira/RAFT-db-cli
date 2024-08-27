package node_test

import (
	"context"
	"errors"
	"testing"
	"time"
	"warehouse/internal/model"
	"warehouse/internal/raft-cluster/node"
	"warehouse/pkg/raft/raft_cluster_v1"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLoadLogNotFollower(t *testing.T) {
	node := &node.ClusterNodeServer{
		IdNode:   1,
		Term:     100,
		State:    node.Lead,
		Logs:     make([]model.Instance, 0),
		SizeLogs: 0,
	}
	logInfo := raft_cluster_v1.LogInfo{
		Id:         uuid.New().String(),
		Term:       100,
		JsonString: "test",
	}

	_, err := node.LoadLog(context.Background(), &logInfo)
	if err == nil || err.Error() != "forbidden, not saved" {
		t.Errorf("expected error 'forbidden, not saved', got %v", err)
	}
}

func TestLoadLogLowerTerm(t *testing.T) {
	node := &node.ClusterNodeServer{
		IdNode:   1,
		Term:     100,
		State:    node.Follower,
		Logs:     make([]model.Instance, 0),
		SizeLogs: 0,
	}
	logInfo := raft_cluster_v1.LogInfo{
		Id:         uuid.New().String(),
		Term:       99,
		JsonString: "test",
	}

	_, err := node.LoadLog(context.Background(), &logInfo)
	if err == nil || err.Error() != "leader is not legitimate, not saved" {
		t.Errorf("expected error 'leader is not legitimate, not saved', got %v", err)
	}
}

func TestLoadLogPanicHandling(t *testing.T) {
	node := &node.ClusterNodeServer{
		IdNode:   1,
		Term:     100,
		State:    node.Follower,
		Logs:     make([]model.Instance, 0),
		SizeLogs: 0,
	}

	// Simulating an error condition by inducing a panic through incorrect index handling
	logInfo := raft_cluster_v1.LogInfo{
		Id:         uuid.New().String(),
		Term:       100,
		JsonString: "test",
		Index:      1, // Out of bounds index to induce a panic
	}

	_, err := node.LoadLog(context.Background(), &logInfo)
	if err == nil || err.Error() != "panic handled on node 1" {
		t.Errorf("expected error 'panic handled on node 1', got %v", err)
	}
}

func TestLoadLogContextTimeout(t *testing.T) {
	node := &node.ClusterNodeServer{
		IdNode:   1,
		Term:     100,
		State:    node.Follower,
		Logs:     make([]model.Instance, 0),
		SizeLogs: 0,
	}
	logInfo := raft_cluster_v1.LogInfo{
		Id:         uuid.New().String(),
		Term:       100,
		JsonString: "test",
	}

	// Creating a context with a very short timeout to trigger context cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(10 * time.Millisecond) // Ensure the context has expired

	_, err := node.LoadLog(ctx, &logInfo)
	if err == nil || !errors.Is(err,  context.Canceled) {
		t.Errorf("expected ctx.Canceled', got %v", err)
	}
}

func TestLoadLogSuccess(t *testing.T) {
	node := &node.ClusterNodeServer{
		IdNode:   1,
		Term:     100,
		State:    node.Follower,
		Logs:     make([]model.Instance, 0),
		SizeLogs: 0,
	}
	logInfo := raft_cluster_v1.LogInfo{
		Id:         uuid.New().String(),
		Term:       100,
		JsonString: "test",
	}

	res, err := node.LoadLog(context.Background(), &logInfo)
	assert.NoError(t, err)
	assert.Equal(t, node.Term, res.Term)
	assert.Equal(t, 1, node.SizeLogs)
	assert.Equal(t, "test", node.Logs[0].Content.Name)
}
