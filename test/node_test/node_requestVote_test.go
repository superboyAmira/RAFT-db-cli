package node_test

import (
	"context"
	"testing"
	"warehouse/internal/model"
	"warehouse/internal/raft-cluster/node"
	"warehouse/pkg/raft/raft_cluster_v1"

	"github.com/google/uuid"
)

func TestRequestVoteContextDone(t *testing.T) {
	node := &node.ClusterNodeServer{
		IdNode:   1,
		Term:     100,
		State:    node.Follower,
		Logs:     []model.Instance{newLogInstance(100, "log1")},
		SizeLogs: 1,
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Сразу отменяем контекст

	req := &raft_cluster_v1.RequestVoteRequest{
		Term:         101,
		LastLogTerm:  100,
		LastLogIndex: 1,
	}

	resp, err := node.RequestVote(ctx, req)
	if err == nil || err.Error() != "[node.go:]ctx done, not saved" {
		t.Errorf("expected error '[node.go:]ctx done, not saved', got %v", err)
	}
	if resp.Vote {
		t.Errorf("expected Vote to be false, got %v", resp.Vote)
	}
}

func TestRequestVoteTermGreater(t *testing.T) {
	node := &node.ClusterNodeServer{
		IdNode:   1,
		Term:     101,
		State:    node.Follower,
		Logs:     []model.Instance{newLogInstance(100, "log1")},
		SizeLogs: 1,
	}

	req := &raft_cluster_v1.RequestVoteRequest{
		Term:         100,
		LastLogTerm:  100,
		LastLogIndex: 1,
	}

	resp, err := node.RequestVote(context.Background(), req)
	if err == nil || err.Error() != "voter's term greater, candidate not legitimate" {
		t.Errorf("expected error 'voter's term greater, candidate not legitimate', got %v", err)
	}
	if resp.Vote {
		t.Errorf("expected Vote to be false, got %v", resp.Vote)
	}
}

func TestRequestVoteChangeStateToFollower(t *testing.T) {
	node1 := &node.ClusterNodeServer{
		IdNode:   1,
		Term:     100,
		State:    node.Lead,
		Logs:     []model.Instance{newLogInstance(100, "log1")},
		SizeLogs: 1,
	}

	req := &raft_cluster_v1.RequestVoteRequest{
		Term:         101,
		LastLogTerm:  100,
		LastLogIndex: 1,
	}

	resp, err := node1.RequestVote(context.Background(), req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !resp.Vote {
		t.Errorf("expected Vote to be true, got %v", resp.Vote)
	}
	if node1.State != node.Follower {
		t.Errorf("expected State to be Follower, got %v", node1.State)
	}
	if node1.Term != 101 {
		t.Errorf("expected Term to be 101, got %v", node1.Term)
	}
}

func TestRequestVoteLastLogTermGreater(t *testing.T) {
	node := &node.ClusterNodeServer{
		IdNode:   1,
		Term:     100,
		State:    node.Follower,
		Logs:     []model.Instance{newLogInstance(100, "log1")},
		SizeLogs: 1,
	}

	req := &raft_cluster_v1.RequestVoteRequest{
		Term:         101,
		LastLogTerm:  101,
		LastLogIndex: 1,
	}

	resp, err := node.RequestVote(context.Background(), req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !resp.Vote {
		t.Errorf("expected Vote to be true, got %v", resp.Vote)
	}
	if node.Term != 101 {
		t.Errorf("expected Term to be 101, got %v", node.Term)
	}
}

func TestRequestVoteLastLogTermSmaller(t *testing.T) {
	node := &node.ClusterNodeServer{
		IdNode:   1,
		Term:     100,
		State:    node.Follower,
		Logs:     []model.Instance{newLogInstance(100, "log1")},
		SizeLogs: 1,
	}

	req := &raft_cluster_v1.RequestVoteRequest{
		Term:         101,
		LastLogTerm:  99,
		LastLogIndex: 1,
	}

	resp, err := node.RequestVote(context.Background(), req)
	if err == nil || err.Error() != "voter's term last log greater, candidate not legitimate" {
		t.Errorf("expected error 'voter's term last log greater, candidate not legitimate', got %v", err)
	}
	if resp.Vote {
		t.Errorf("expected Vote to be false, got %v", resp.Vote)
	}
}

func TestRequestVoteLastLogIndexGreater(t *testing.T) {
	node := &node.ClusterNodeServer{
		IdNode:   1,
		Term:     100,
		State:    node.Follower,
		Logs:     []model.Instance{newLogInstance(100, "log1")},
		SizeLogs: 1,
	}

	req := &raft_cluster_v1.RequestVoteRequest{
		Term:         101,
		LastLogTerm:  100,
		LastLogIndex: 2,
	}

	resp, err := node.RequestVote(context.Background(), req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !resp.Vote {
		t.Errorf("expected Vote to be true, got %v", resp.Vote)
	}
	if node.Term != 101 {
		t.Errorf("expected Term to be 101, got %v", node.Term)
	}
}

func TestRequestVoteLastLogIndexSmaller(t *testing.T) {
	node := &node.ClusterNodeServer{
		IdNode:   1,
		Term:     100,
		State:    node.Follower,
		Logs:     []model.Instance{newLogInstance(100, "log1")},
		SizeLogs: 1,
	}

	req := &raft_cluster_v1.RequestVoteRequest{
		Term:         101,
		LastLogTerm:  100,
		LastLogIndex: 0,
	}

	resp, err := node.RequestVote(context.Background(), req)
	if err == nil || err.Error() != "voter's logs more complete, candidate not legitimate" {
		t.Errorf("expected error 'voter's logs more complete, candidate not legitimate', got %v", err)
	}
	if resp.Vote {
		t.Errorf("expected Vote to be false, got %v", resp.Vote)
	}
}

func newLogInstance(term int, data string) model.Instance {
	return model.Instance{
		Id:      uuid.New(),
		Content: model.JsonData{Name: data},
		Term:    int64(term),
	}
}
