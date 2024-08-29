/*
This is a common implementation of the DRAFT protocol node for distributed fault-tolerant storage.
Complies with the principles of the RAFT protocol:
1. Log synchronization
2. Replication and Commit
3. Working with HeartBeat
4. Thread Safety
5. Node state management and election logic
*/
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

// Global cluster settings
type ClusterSettings struct {
	Host string `json:"cluster_host"`
	Port string

	Size              int `json:"cluster_size"`
	Quorum            int `json:"quorum"` // roundUp(initClusterSize)/2
	ReplicationFactor int `json:"replication_factor"`

	HeartBeatTimeout        time.Duration // follower mutate to candidate after the HeartBeatTimeout expires
	HeartBeatIntervalLeader time.Duration // lead must send a heartbeats to followers every second
}

// const (
// 	initClusterSize   = 3
// 	quorum            = 2
// 	replicationFactor = 2

//
// 	HeartBeatTimeout = time.Second * 3
//
// 	HeartBeatIntervalLeader = time.Second * 1
// )

// Node current role in custer
type StateType int

const (
	Follower StateType = iota
	Candidate
	Lead
)

type ClusterNodeServer struct {
	Settings *ClusterSettings

	IdNode   int
	Logs     []model.Instance
	SizeLogs int

	Term  int64
	State StateType
	// pointers to all nodes (poiteres, beacuse we are using mutex)
	Network            []*ClusterNodeServer
	LeadId             int
	timeLastHeartBreak map[int]time.Time

	mu sync.RWMutex

	raft_cluster_v1.UnimplementedClusterNodeServer
}

// constructor
func New(idNode int, state StateType, leadId int, network []*ClusterNodeServer, sett *ClusterSettings) *ClusterNodeServer {
	return &ClusterNodeServer{
		Settings:           sett,
		IdNode:             idNode,
		Logs:               make([]model.Instance, 0), // Инициализация пустого слайса
		SizeLogs:           0,
		Term:               0,
		State:              state,
		Network:            network,
		LeadId:             leadId,
		timeLastHeartBreak: make(map[int]time.Time), // Инициализация пустой карты
		mu:                 sync.RWMutex{},
	}
}

// Start node server on random free port and fixed host form cfg
func (r *ClusterNodeServer) Serve(stop <-chan struct{}) error {
	// listen, err := net.Listen("tcp", config.Port)

	// if err != nil {
	// 	log.Fatalf("Failed to listen: %v", err)
	// }

	// server := grpc.NewServer()
	// raft_cluster_v1.RegisterClusterNodeServer(server, &transmitter.ConnServer{})

	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT, syscall.SIGTSTP)

	// log.Printf("Server is running on port %s", config.Port)
	// go func() {
	// 	if err := server.Serve(listen); err != nil {
	// 		log.Fatalf("Failed to serve: %v", err)
	// 	}
	// }()
	<-stop
	// sercer.GracefulyStop()

	return nil
}

func (r *ClusterNodeServer) Stop() {}

// Writing data to the node storage
func (r *ClusterNodeServer) LoadLog(ctx context.Context, req *raft_cluster_v1.LogInfo) (*raft_cluster_v1.LogAccept, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	// check node role
	if r.State == Candidate {
		return &raft_cluster_v1.LogAccept{Term: r.Term}, errors.New("forbidden, not saved")
	}
	if r.Term > req.Term {
		return &raft_cluster_v1.LogAccept{Term: r.Term}, errors.New("leader is not legitimate, not saved")
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
		return &raft_cluster_v1.LogAccept{Term: r.Term}, context.Canceled

	case err := <-commit:
		if err != nil {
			return &raft_cluster_v1.LogAccept{Term: r.Term}, err
		} else {
			r.Term = req.Term
			r.SizeLogs = len(r.Logs)
			return &raft_cluster_v1.LogAccept{Term: r.Term}, nil
		}
	}
}

// Requesting vote from Candidate to Follower
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

		if r.Term < req.Term {
			if r.State == Lead || r.State == Candidate {
				r.State = Follower
			}
			commit <- nil
		} else if r.Term > req.Term {
			commit <- errors.New("voter's term greater, candidate not legitimate")
		} else {
			// if request was sendet to high level nodes, they changing their state
			// it doesn`t mean that they always responsed 'yes', need check relevance logs

			// check relevance logs
			if req.LastLogTerm > r.Logs[r.SizeLogs-1].Term {
				// term more than our -> relevated, vote yes
				commit <- nil
			} else if req.LastLogTerm < r.Logs[r.SizeLogs-1].Term {
				commit <- errors.New("voter's term last log greater, candidate not legitimate")
			} else {
				// if terms equal checking lenght of logs
				switch req.LastLogIndex > int64(r.SizeLogs)-1 {
				case true:
					if r.State == Lead || r.State == Candidate {
						r.State = Follower
					}
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
		return &raft_cluster_v1.RequestVoteResponse{Term: r.Term}, context.Canceled
	case err := <-commit:
		if err != nil {
			return &raft_cluster_v1.RequestVoteResponse{Term: r.Term}, err
		} else {
			r.mu.Lock()
			r.Term = req.Term
			r.mu.Unlock()
			return &raft_cluster_v1.RequestVoteResponse{Term: r.Term}, nil
		}
	}
}

// after the HeartBeatTimeout expires folower start election and send RequestVote to all nodes [Candidate method]
// TODO: HeartBeat timer to new lead
func (r *ClusterNodeServer) StartElection(ctx context.Context, req *raft_cluster_v1.Empty) (*raft_cluster_v1.ElectionDecision, error) {
	if r.State != Follower {
		return &raft_cluster_v1.ElectionDecision{Term: r.Term}, errors.New("state unsuppoted to election")
	}

	commit := make(chan error, 1)
	go func() {
		select {
		case <-ctx.Done():
			return
		default:
		}

		r.mu.Lock()
		defer r.mu.Unlock()
		// we update r.Term with new election
		r.Term++
		LastLogTerm := r.Term
		LastLogIndex := 0
		if r.SizeLogs > 0 {
			LastLogTerm = r.Logs[r.SizeLogs-1].Term
			LastLogIndex = r.SizeLogs - 1
		}
		// mutate to Candidate
		r.State = Candidate
		r.LeadId = -1

		ballotbox := 0
		for _, node := range r.Network {
			bulletin, err := node.RequestVote(ctx, &raft_cluster_v1.RequestVoteRequest{
				Term:         r.Term,
				LastLogTerm:  LastLogTerm,
				LastLogIndex: int64(LastLogIndex),
			})

			if err != nil {
				// network err
				if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
					continue
				}
				// update our node data and state, if we are not legigimate
				r.Term = bulletin.Term
				r.State = Follower
				r.Logs = node.Logs // update
				r.LeadId = node.IdNode
				commit <- err
				return
			} else {
				ballotbox++
			}
		}

		// checking qourum requirement
		// ballotbox + candidate votes to yourself
		if ballotbox+1 >= r.Settings.Quorum {
			r.State = Lead
			// send all nodes, that current node became a lead
			for _, node := range r.Network {
				_, err := node.SetLeader(ctx, &raft_cluster_v1.LeadInfo{IdLeader: int64(r.IdNode), Term: r.Term})
				// if a node has more up-to-date information,
				// we send a request to nominate its candidacy in the elections
				if err != nil {
					// network err with retry
					if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
						time.Sleep(time.Millisecond * 10)
						_, err = node.SetLeader(ctx, &raft_cluster_v1.LeadInfo{IdLeader: int64(r.IdNode)})
						if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
							continue
						}
					}
					r.State = Follower
					commit <- err
					return
				}
			}
			commit <- nil
		} else {
			r.State = Follower
			commit <- errors.New("qourum not satisfied")
		}
	}()

	// control for context
	select {
	case <-ctx.Done():
		return &raft_cluster_v1.ElectionDecision{Term: r.Term}, context.Canceled
	case err := <-commit:
		return &raft_cluster_v1.ElectionDecision{Term: r.Term}, err
	}
}

// Follower
// WARN:without Start Election, it requested from candidates StartElection, if needed
func (r *ClusterNodeServer) SetLeader(ctx context.Context, req *raft_cluster_v1.LeadInfo) (*raft_cluster_v1.LeadAccept, error) {
	if r.Term > req.Term {
		return &raft_cluster_v1.LeadAccept{
			Term: r.Term,
		}, errors.New("lead isn't legitimate")
	}

	commit := make(chan error, 1)
	go func() {
		r.mu.Lock()
		defer r.mu.Unlock()

		r.LeadId = int(req.IdLeader)
		r.Term = req.Term
		commit <- nil
	}()

	// control for context
	select {
	case <-ctx.Done():
		return &raft_cluster_v1.LeadAccept{Term: r.Term}, context.Canceled
	case err := <-commit:
		return &raft_cluster_v1.LeadAccept{Term: r.Term}, err
	}
}

// Lead
func (r *ClusterNodeServer) SendHeartBeat(ctx context.Context, req *raft_cluster_v1.Empty) (*raft_cluster_v1.Empty, error) {
	if r.State == Lead {
		ticker := time.NewTicker(r.Settings.HeartBeatIntervalLeader)
		for {
			select {
			case <-ctx.Done():
				return &raft_cluster_v1.Empty{}, context.Canceled
			case <-ticker.C:
				r.mu.RLock()
				prevLogTerm := r.Term
				if r.SizeLogs > 0 {
					prevLogTerm = r.Logs[r.SizeLogs-1].Term
				}
				heartbeat := &raft_cluster_v1.HeartBeatRequest{
					Term:         r.Term,
					LeaderId:     int64(r.IdNode),
					PrevLogIndex: int64(r.SizeLogs) - 1,
					PrevLogTerm:  prevLogTerm,
				}
				r.mu.RUnlock()

				var wg sync.WaitGroup

				legitimate := true
				for _, follower := range r.Network {
					wg.Add(1)
					go func(node *ClusterNodeServer) {
						defer wg.Done()
						_, err := node.ReciveHeartBeat(ctx, heartbeat)
						if err != nil {
							r.mu.Lock()
							if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
								r.Network = slices.Delete(r.Network, node.IdNode-1, node.IdNode)
							} else {
								legitimate = false
							}
							r.mu.Unlock()
						}
					}(follower)
					if !legitimate {
						r.mu.Lock()
						r.State = Follower
						r.mu.Unlock()
						return &raft_cluster_v1.Empty{}, errors.New("lead not legitimate")
					}
				}
				wg.Wait()
			}
		}
	} else {
		return &raft_cluster_v1.Empty{}, errors.ErrUnsupported
	}
}

// TODO: check request(lead), maybe restarted node
func (r *ClusterNodeServer) ReciveHeartBeat(ctx context.Context, req *raft_cluster_v1.HeartBeatRequest) (*raft_cluster_v1.HeartBeatResponse, error) {
	r.timeLastHeartBreak[int(req.LeaderId)] = time.Now()

	commit := make(chan error, 1)
	go func() {
		r.mu.RLock()
		defer r.mu.RUnlock()
		// check relevance leader
		if r.SizeLogs == 0 {
			commit <- nil
		} else if req.PrevLogTerm > r.Logs[r.SizeLogs-1].Term {
			commit <- nil
		} else if req.PrevLogTerm < r.Logs[r.SizeLogs-1].Term {
			commit <- errors.New("node's term last log greater, leader not legitimate")
		} else {
			// if terms equal checking lenght of logs
			switch req.PrevLogIndex > int64(r.SizeLogs)-1 {
			case true:
				if r.State == Lead || r.State == Candidate {
					r.State = Follower
				}
				commit <- nil
			case false:
				commit <- errors.New("node's logs more complete, leader not legitimate")
			}
		}
	}()

	// control for context
	select {
	case <-ctx.Done():
		return &raft_cluster_v1.HeartBeatResponse{Term: r.Term}, context.Canceled
	case err := <-commit:
		if err != nil {
			return &raft_cluster_v1.HeartBeatResponse{Term: r.Term}, err
		} else {
			r.mu.Lock()
			r.Term = req.Term
			r.LeadId = int(req.LeaderId)
			r.Logs = r.Network[r.LeadId].Logs
			r.mu.Unlock()
			return &raft_cluster_v1.HeartBeatResponse{Term: r.Term}, nil
		}
	}
}

// Load log to Cluster. This is an abstract method for changing data in the entire cluster.
func (r *ClusterNodeServer) Append(ctx context.Context, req *raft_cluster_v1.LogLeadRequest) (*raft_cluster_v1.Empty, error) {
	if r.State != Lead {
		_, err := r.Network[r.LeadId].Append(ctx, req)
		return &raft_cluster_v1.Empty{}, err
	} else {
		if req == nil {
			return &raft_cluster_v1.Empty{}, errors.New("nil reference req")
		}
		log := &raft_cluster_v1.LogInfo{
			Id:         req.Id,
			Term:       r.Term,
			Index:      int64(r.SizeLogs),
			JsonString: req.JsonString,
		}
		_, err := r.LoadLog(ctx, log)
		if err != nil {
			return &raft_cluster_v1.Empty{}, err
		}

		// replication process
		loaded := 1
		for _, node := range r.Network {
			if loaded == r.Settings.Quorum {
				break
			}
			resp, err := node.LoadLog(ctx, log)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
					r.Network = slices.Delete(r.Network, node.IdNode-1, node.IdNode)
					continue
				} else {
					// TODO: think about Election
					if resp.Term > r.Term {
						r.State = Follower
					}
					r.Logs = r.Logs[:len(r.Logs)-1]
					// rollback our log form prev nodes.
					// ... TODO:
					return &raft_cluster_v1.Empty{}, err
				}
			}
			loaded++
		}
		if loaded < r.Settings.Quorum {
			return &raft_cluster_v1.Empty{}, errors.New("quorum not required")
		} else {
			return &raft_cluster_v1.Empty{}, nil
		}
	}
}
