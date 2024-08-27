package node

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"sync"
	"time"
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

	// follower mutate to candidate after the HeartBeatTimeout expires 
	HeartBeatTimeout = time.Second * 3
	// lead must send a heartbeats to followers every second
	HeartBeatIntervalLeader = time.Second * 1
)

// Node current role in custer
type StateType int

const (
	Follower StateType = iota
	Candidate
	Lead
)

type ClusterNodeServer struct {
	IdNode int
	Logs     []model.Instance
	mu       sync.RWMutex
	SizeLogs int

	Term   int64
	State  StateType
	Network []ClusterNodeServer
	LeadId int

	raft_cluster_v1.UnimplementedClusterNodeServer
}

// Writing data to the node storage [Follower method]
func (r *ClusterNodeServer) LoadLog(ctx context.Context, req *raft_cluster_v1.LogInfo) (*raft_cluster_v1.LogAccept, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// check node role
	if r.State != Follower {
		return &raft_cluster_v1.LogAccept{ Term: r.Term }, errors.New("forbidden, not saved")
	}
	if r.Term > req.Term {
		return &raft_cluster_v1.LogAccept{ Term: r.Term }, errors.New("leader is not legitimate, not saved")
	}

	// if role is follower -> saving
	commit := make(chan error, 1) // buffer chan for correct select construction
	go func() {
		// for possible routine leaks with ctx cancel
		select {
		case <-ctx.Done():
			return
		default:
		}

		defer func() {
			// rollback
			if rec := recover(); rec != nil {
				// if log added delete him
				if r.SizeLogs > 0 {
					if r.Logs[r.SizeLogs-1].Id.String() == req.Id {
						r.Logs = r.Logs[:len(r.Logs)-1]
					}
				}
				// DEBUG needed!!!
				// if term is correct, deactivate node
				// if req.Term > r.term {
				// }
				commit <- errors.New("panic handled on node " + strconv.FormatInt(int64(r.IdNode), 10))
			}
		}()

		r.Logs = slices.Insert(r.Logs, int(req.Index), model.Instance{
			Id:      uuid.New(),
			Content: model.JsonData{Name: req.JsonString},
			Term:    req.Term,
		})
		commit <- nil
	}()

	// control for context
	select {
	case <-ctx.Done():
		if r.SizeLogs > 0 {
			if r.Logs[r.SizeLogs-1].Id.String() == req.Id {
				r.Logs = slices.Delete(r.Logs, int(r.SizeLogs-2), int(r.SizeLogs-1))
			}
		}
		return &raft_cluster_v1.LogAccept{ Term: r.Term }, errors.New("ctx done, not saved")

	case err := <-commit:
		if err != nil {
			return &raft_cluster_v1.LogAccept{ Term: r.Term }, err
		} else {
			r.mu.Lock()
			r.Term = req.Term
			r.SizeLogs = len(r.Logs)
			r.mu.Unlock()
			return &raft_cluster_v1.LogAccept{ Term: r.Term }, nil
		}
	}
}

// Requesting vote from Candidate to Follower [Follower method]
func (r *ClusterNodeServer) RequestVote(ctx context.Context, req *raft_cluster_v1.RequestVoteRequest) (*raft_cluster_v1.RequestVoteResponse, error) {
	commit := make(chan error, 1)
	go func() {
		// for possible routine leaks with ctx cancel
		select {
		case <-ctx.Done():
			return
		default:
		}

		r.mu.Lock()
		defer r.mu.Unlock()

		if r.Term > req.Term {
			commit <- errors.New("voter's term greater, candidate not legitimate")
		} else {
			// if request was sendet to high level nodes, they changing their state
			// it doesn`t mean that they always responsed 'yes', need check relevance logs
			if r.State == Lead || r.State == Candidate {
				r.State = Follower
			}
			// check relevance logs
			if req.LastLogTerm > r.Logs[r.SizeLogs-1].Term {
				// term more than our -> relevated, vote yes
				commit <- nil
			} else if req.LastLogTerm < r.Logs[r.SizeLogs-1].Term {
				commit <- errors.New("voter's term last log greater, candidate not legitimate")
			} else {
				// if terms equal checking lenght of logs
				switch req.LastLogIndex >= int64(r.SizeLogs) {
				case true:
					commit <- nil
				case false:
					commit <- errors.New("voter's logs more complete, candidate not legitimate")
				}
			}
		}
	}()

	// control for context
	select {
	case <-ctx.Done():
		return &raft_cluster_v1.RequestVoteResponse{ Term: r.Term }, errors.New("ctx done, not saved")
	case err := <-commit:
		if err != nil {
			return &raft_cluster_v1.RequestVoteResponse{ Term: r.Term }, err
		} else {
			r.mu.Lock()
			r.Term = req.Term
			r.mu.Unlock()
			return &raft_cluster_v1.RequestVoteResponse{ Term: r.Term }, nil
		}
	}
}

// after the HeartBeatTimeout expires folower start election and send RequestVote to all nodes [Candidate method]
func (r *ClusterNodeServer) StartElection(ctx context.Context, req *raft_cluster_v1.Empty) (*raft_cluster_v1.ElectionDecision, error) {
	if r.State != Follower {
		return &raft_cluster_v1.ElectionDecision{ Term: r.Term }, errors.New("state unsuppoted to election")
	}

	commit := make(chan error, 1)
	go func ()  {
		select {
		case <-ctx.Done():
			return
		default:
		}

		// we update r.Term with new election
		LastLogTerm := r.Term + 1
		LastLogIndex := 0
		if r.SizeLogs > 0 {
			LastLogTerm = r.Logs[r.SizeLogs-1].Term
			LastLogIndex = r.SizeLogs-1
		}
		// mutate to Candidate
		r.State = Candidate

		for _, node := range r.Network {
			bulletin, err := node.RequestVote(ctx, &raft_cluster_v1.RequestVoteRequest{
				Term: r.Term,
				LastLogTerm: LastLogTerm,
				LastLogIndex: int64(LastLogIndex),
			})

			if err != nil {
				if bulletin.Term > r.Term {
					r.Term = bulletin.Term
				}
				r.State = Follower
				commit <- err
			}
		}
	}()


	// control for context
	select {
	case <-ctx.Done():
		return &raft_cluster_v1.ElectionDecision{ Term: r.Term }, errors.New("ctx done, not saved")
	case err := <-commit:
		return &raft_cluster_v1.ElectionDecision{ Term: r.Term }, err
	}
}

// // Follower
// SetLeader(context.Context, *LeadInfo) (*LeadAccept, error)
// // Lead
// SendHeartBeat(context.Context, *HeartBeatRequest) (*HeartBeatResponse, error)