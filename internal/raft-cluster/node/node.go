package node

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"warehouse/internal/model"
	"warehouse/pkg/raft/raft_cluster_v1"

	"github.com/google/uuid"
)

/*
Как только кластер получает лидера, он может принимать новые записи журнала.
Клиент может запросить лидера добавить новую запись журнала,
которая представляет собой непрозрачный двоичный объект в Raft.
Затем лидер записывает запись в долговременное хранилище и пытается реплицировать ее
в кворум последователей. Как только запись журнала считается зафиксированной ,
ее можно применить к конечному автомату.
Конечный автомат зависит от приложения и реализуется с помощью интерфейса.
*/

// Global cluster settings
const (
	initClusterSize   = 3
	quorum            = 2 // roundUp(initClusterSize)/2
	replicationFactor = 2
)

// Node current role in custer
type StateType int

const (
	Follower StateType = iota
	Candidate
	Lead
)

type ClusterNodeServer struct {
	idNode int
	term   int64
	state  StateType

	logs     []model.Instance
	sizeLogs int

	raft_cluster_v1.UnimplementedClusterNodeServer
}

// Writing data to the node storage [Follower method]
func (r *ClusterNodeServer) LoadLog(ctx context.Context, req *raft_cluster_v1.LogInfo) (*raft_cluster_v1.LogAccept, error) {
	Saved := &raft_cluster_v1.LogAccept{
		Accept: true,
		Term:   r.term,
	}
	NotSaved := &raft_cluster_v1.LogAccept{
		Accept: false,
		Term:   r.term,
	}

	// check node role
	if r.state != Follower {
		return NotSaved, errors.New("forbidden, not saved")
	}
	if r.term > req.Term {
		return NotSaved, errors.New("leader is not legitimate, not saved")
	}

	// if role is follower -> saving
	commit := make(chan error)
	go func() {
		defer func() {
			// rollback
			if rec := recover(); rec != nil {
				// if log added delete him
				if r.logs[r.sizeLogs-1].Id.String() == req.Id {
					r.logs = slices.Delete(r.logs, int(r.sizeLogs-2), int(r.sizeLogs-1))
				}
				// DEBUG needed!!!
				// if term is correct, deactivate node
				// if req.Term > r.term {
				// }
				commit <- errors.New("panic handled on node " + strconv.FormatInt(int64(r.idNode), 10))
			}
		}()

		r.logs = slices.Insert(r.logs, int(req.Index), model.Instance{
			Id:      uuid.New(),
			Content: model.JsonData{Name: req.JsonString},
			Term:    req.Term,
		})
		commit <- nil
	}()

	// control for context
	select {
	case <-ctx.Done():
		if r.logs[r.sizeLogs-1].Id.String() != req.Id {
			return NotSaved, errors.New("ctx done, not saved")
		}
		return Saved, nil
	case err := <-commit:
		if err != nil {
			return NotSaved, err
		} else {
			r.term = req.Term
			r.sizeLogs = len(r.logs)
			return Saved, nil
		}
	}
}

// Requesting vote from Candidate to Follower [Follower method]
func (r *ClusterNodeServer) RequestVote(ctx context.Context, req *raft_cluster_v1.RequestVoteRequest) (*raft_cluster_v1.RequestVoteResponse, error) {
	yes := &raft_cluster_v1.RequestVoteResponse{
		Vote: true,
		Term: r.term,
	}
	no := &raft_cluster_v1.RequestVoteResponse{
		Vote: false,
		Term: r.term,
	}
	commit := make(chan error)
	go func() {
		defer func() {
			// panic handler
			if rec := recover(); rec != nil {
				commit <- errors.New("panic handled on node " + strconv.FormatInt(int64(r.idNode), 10))
			}
		}()
		if r.term > req.Term {
			commit <- errors.New("voter's term greather, candidate not legitimate")
		} else {
			// if request was sendet to high level nodes, they changing their state
			// it doesn`t mean that they always responsed 'yes', need check relevance logs
			if r.state == Lead || r.state == Candidate {
				r.state = Follower
			}
			// check relevance logs
			if req.LastLogTerm > r.logs[r.sizeLogs-1].Term {
				// term more than our -> relevated, vote yes
				commit <- nil
			} else if req.LastLogTerm < r.logs[r.sizeLogs-1].Term {
				commit <- errors.New("voter's term last log greather, candidate not legitimate")
			} else {
				// if terms equal checking lenght of logs
				switch req.LastLogIndex > int64(r.sizeLogs) {
				case true:
					commit <- nil
				case false:
					commit <- errors.New("voters's logs more completely, candidate not legitimate")
				}
			}
		}
	}()

	// control for context
	select {
	case <-ctx.Done():
		return no, errors.New("[node.go:]ctx done, not saved")
	case err := <-commit:
		if err != nil {
			return no, err
		} else {
			r.term = req.Term
			return yes, nil
		}
	}
}

// // Follower
// SetLeader(context.Context, *LeadInfo) (*LeadAccept, error)
// // Lead
// SendHeartBeat(context.Context, *HeartBeatRequest) (*HeartBeatResponse, error)
// // Candidate
// StartElection(context.Context, *Empty) (*ElectionDecision, error)
