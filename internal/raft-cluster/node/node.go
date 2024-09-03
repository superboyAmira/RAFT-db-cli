/*
This is a common implementation of the RAFT protocol node for distributed fault-tolerant storage.
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
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
	"warehouse/internal/model"
	"warehouse/pkg/raft/raft_cluster_v1"

	"github.com/google/uuid"
	"google.golang.org/grpc"
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
	ElectionTimeout         time.Duration
}

var Log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

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

	Network       []raft_cluster_v1.ClusterNodeClient
	LeadId        int
	electionTimer *time.Timer

	mu         sync.RWMutex
	nodeCtx    context.Context
	nodeCancel context.CancelFunc

	raft_cluster_v1.UnimplementedClusterNodeServer
}

/*
----------------Start Cluster-------------------
*/

// Cluster node constructor
func New(idNode int, state StateType, leadId int, term int, sett *ClusterSettings) *ClusterNodeServer {
	t := &ClusterNodeServer{
		Settings: sett,
		IdNode:   idNode,
		Logs:     make([]model.Instance, 0),
		SizeLogs: 0,
		Term:     int64(term),
		State:    state,
		Network:  make([]raft_cluster_v1.ClusterNodeClient, 0),
		LeadId:   leadId,
		mu:       sync.RWMutex{},
	}
	return t
}

// Start node server on random free port and fixed host form cfg
func (r *ClusterNodeServer) Serve(ctx context.Context) error {
	// copy to local context
	// this useful to untie the execution of a node from the context of the entire cluster.
	// When using a shared cluster, when they fail, the context of the entire cluster will be canceled.
	r.nodeCtx, r.nodeCancel = context.WithCancel(ctx)

	listen, err := net.Listen("tcp", r.Settings.Port)
	if err != nil {
		Log.Warn("Failed to listen", slog.Int("node_id", r.IdNode), slog.String("error", err.Error()))
		return err
	}
	server := grpc.NewServer()
	raft_cluster_v1.RegisterClusterNodeServer(server, r)

	Log.Info("Node is running on port", slog.Int("node_id", r.IdNode), slog.String("address", listen.Addr().String()))

	// start server
	go func() {
		if err := server.Serve(listen); err != nil {
			Log.Error("Failed to serve ", slog.Int("node_id", r.IdNode), slog.String("error", err.Error()))
		}
	}()

	time.Sleep(r.Settings.ElectionTimeout * 10)

	r.BecameFollower(r.nodeCtx)
	// waiting cancel() from parent context in manager
	<-r.nodeCtx.Done()
	if r.SizeLogs > 0 {
		Log.Info("Death note about node", slog.Int("node_id", r.IdNode), slog.Int("leadID", r.LeadId), slog.Int("term", int(r.Term)), slog.String("last_content", r.Logs[r.SizeLogs-1].Content.Name))
	} else {
		Log.Info("Death note about node", slog.Int("node_id", r.IdNode), slog.Int("leadID", r.LeadId), slog.Int("term", int(r.Term)), slog.String("address", listen.Addr().String()))
	}
	Log.Info("Shutdown...", slog.Int("node_id", r.IdNode), slog.String("address", listen.Addr().String()))
	server.GracefulStop()
	Log.Info("Stopped Gracefully...", slog.Int("node_id", r.IdNode), slog.String("address", listen.Addr().String()))
	return nil
}

/*
----------------State Mutations-------------------
*/

func (r *ClusterNodeServer) BecameLead(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
	}
	Log.Debug("Became Lead", slog.Int("node_id", r.IdNode), slog.String("address", r.Settings.Port))
	r.mu.Lock()
	r.State = Lead
	r.LeadId = r.IdNode
	r.mu.Unlock()
	go r.HeartBeatTicker(r.nodeCtx)
}

func (r *ClusterNodeServer) BecameFollower(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
	}
	Log.Debug("Became Folower", slog.Int("node_id", r.IdNode), slog.String("address", r.Settings.Port))
	r.mu.Lock()
	r.State = Follower
	// for follower to lead mutation
	r.mu.Unlock()
	r.ResetElectionTimer(r.nodeCtx)
}

func (r *ClusterNodeServer) BecameCandidate(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	Log.Debug("Became Candidate", slog.Int("nodeID", r.IdNode), slog.String("address", r.Settings.Port))
	r.mu.Lock()
	// if !r.electionTimer.Stop() {
	// <-r.electionTimer.C
	// }
	r.State = Candidate
	r.mu.Unlock()
	_, err := r.StartElection(r.nodeCtx, &raft_cluster_v1.Empty{})
	if err != nil {
		Log.Debug("Election result", slog.String("port", r.Settings.Port), slog.String("err", err.Error()))
		r.BecameFollower(r.nodeCtx)
		return nil
	}
	r.BecameLead(r.nodeCtx)
	return nil
}

/*
----------------Heartbeat-------------------
*/

func (r *ClusterNodeServer) SetElectionTimeout(ctx context.Context, req *raft_cluster_v1.Empty) (*raft_cluster_v1.Empty, error) {
	r.ResetElectionTimer(r.nodeCtx)
	return nil, nil
}

func (r *ClusterNodeServer) ResetElectionTimer(ctx context.Context) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.electionTimer != nil {
		if !r.electionTimer.Stop() {
			select {
			case <-r.electionTimer.C:
			default:
			}
		}
	}

	// this magic number '4' needed in order to eliminate parallel elections between two nodes
	duration := time.Duration(r.Settings.HeartBeatTimeout + time.Millisecond*time.Duration(r.IdNode*rand.Intn(7)+7))
	r.electionTimer = time.AfterFunc(duration, func() {
		select {
		case <-ctx.Done():
			Log.Error("Election timer cancelled due to context cancellation")
			return
		default:
			Log.Debug("TO CAND", slog.Int("nodeID", r.IdNode), slog.String("address", r.Settings.Port))
			for _, node := range r.Network {
				node.SetElectionTimeout(r.nodeCtx, &raft_cluster_v1.Empty{})
			}
			r.BecameCandidate(r.nodeCtx)
		}
	})

	Log.Debug("EL RESET", "dur", duration, "port", r.Settings.Port)
}

func (r *ClusterNodeServer) HeartBeatTicker(ctx context.Context) {
	Log.Debug("Lead Ticker started", slog.String("port", r.Settings.Port))
	ticker := time.NewTicker(r.Settings.HeartBeatIntervalLeader)
	for {
		select {
		case <-ctx.Done():
			Log.Debug("Lead Disconnected!")
			return
		case <-ticker.C:
			// r.mu.RLock()
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
			// r.mu.RUnlock()

			var wg sync.WaitGroup

			legitimate := true
			for id, follower := range r.Network {
				wg.Add(1)
				go func(client raft_cluster_v1.ClusterNodeClient, id int) {
					defer wg.Done()
					_, err := client.ReciveHeartBeat(r.nodeCtx, heartbeat)
					if err != nil {
						r.mu.Lock()
						if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
							r.Network = slices.Delete(r.Network, id-1, id)
						} else {
							legitimate = false
						}
						r.mu.Unlock()
					}
				}(follower, id)
				if !legitimate {
					wg.Wait()
					r.BecameFollower(r.nodeCtx)
					Log.Debug("Lead Ticker err ended", slog.String("port", r.Settings.Port))
					return
				}
			}
			wg.Wait()
		}
	}
}

// TODO: check request(lead), maybe restarted node
func (r *ClusterNodeServer) ReciveHeartBeat(ctx context.Context, req *raft_cluster_v1.HeartBeatRequest) (*raft_cluster_v1.HeartBeatResponse, error) {
	select {
	case <-ctx.Done():
		return &raft_cluster_v1.HeartBeatResponse{Term: r.Term}, context.Canceled
	default:
		if r.State == Lead || r.State == Candidate {
			r.BecameFollower(r.nodeCtx)
		}
		r.ResetElectionTimer(r.nodeCtx)
		r.mu.Lock()
		r.Term = req.Term
		r.LeadId = int(req.LeaderId)
		r.mu.Unlock()
		Log.Debug("HB Recieved form ", slog.Int("nodeID", r.IdNode), slog.Int("senderID", int(req.LeaderId)))
		return &raft_cluster_v1.HeartBeatResponse{Term: r.Term}, nil
	}
}

/*
----------------Election-------------------
*/

// after the HeartBeatTimeout expires folower start election and send RequestVote to all nodes [Candidate method]
// TODO: HeartBeat timer to new lead
func (r *ClusterNodeServer) StartElection(ctx context.Context, req *raft_cluster_v1.Empty) (*raft_cluster_v1.ElectionDecision, error) {
	Log.Debug("START Election", slog.Int("nodeID", r.IdNode), slog.String("address", r.Settings.Port))
	select {
	case <-ctx.Done():
		return &raft_cluster_v1.ElectionDecision{Term: r.Term}, ctx.Err()
	default:
	}

	r.mu.Lock()
	// we update r.Term with new election
	r.Term++
	r.LeadId = -1
	r.mu.Unlock()

	r.mu.RLock()
	LastLogTerm := r.Term
	LastLogIndex := 0
	if r.SizeLogs > 0 {
		LastLogTerm = r.Logs[r.SizeLogs-1].Term
		LastLogIndex = r.SizeLogs - 1
	}
	r.mu.RUnlock()

	ballotbox := 1
	for id, network_client := range r.Network {
		Log.Debug("Send VoteReq", slog.Int("nodeID", r.IdNode), slog.String("address", r.Settings.Port))
		bulletin, err := network_client.RequestVote(r.nodeCtx, &raft_cluster_v1.RequestVoteRequest{
			Term:         r.Term,
			LastLogTerm:  LastLogTerm,
			LastLogIndex: int64(LastLogIndex),
			SenderId:     int64(r.IdNode),
		})

		if err != nil {
			// network err
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				continue
			}
			if strings.Contains(err.Error(), "refused") {
				continue
			}
			Log.Warn(err.Error())
			// update our node data and state, if we are not legigimate
			r.Term = bulletin.Term
			r.LeadId = id
			// r.BecameFollower(ctx) this func triggered in became candidate in all cases
			return &raft_cluster_v1.ElectionDecision{Term: r.Term}, err
		} else {
			ballotbox++
		}
	}

	Log.Debug("TEMP Election temp results", slog.Int("nodeID", r.IdNode), slog.Int("ballotbox", ballotbox))
	// checking qourum requirement
	// ballotbox + candidate votes to yourself
	if ballotbox >= r.Settings.Quorum {
		r.State = Lead
		logs := make([]*raft_cluster_v1.LogInfo, 0)
		for i, log := range r.Logs {
			logs = append(logs, &raft_cluster_v1.LogInfo{
				Id:         log.Id.String(),
				Term:       log.Term,
				Index:      int64(i),
				JsonString: log.Content.Name,
			})
		}
		// send all nodes, that current node became a lead
		for _, node := range r.Network {
			node.SetLeader(ctx, &raft_cluster_v1.LeadInfo{IdLeader: int64(r.IdNode), Term: r.Term, Logs: logs})
			// _, err := node.SetLeader(ctx, &raft_cluster_v1.LeadInfo{IdLeader: int64(r.IdNode), Term: r.Term})
			// start Lead HeartBeat
			// go r.HeartBeatTicker(ctx)
			// if a node has more up-to-date information,
			// we send a request to nominate its candidacy in the elections
			// if err != nil {
			// 	// network err with retry
			// 	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			// 		time.Sleep(time.Millisecond * 10)
			// 		_, err = node.SetLeader(ctx, &raft_cluster_v1.LeadInfo{IdLeader: int64(r.IdNode)})
			// 		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			// 			continue
			// 		}
			// 	}
			// 	// r.BecameFollower(ctx) this func triggered in became candidate in all cases
			// 	commit <- err
			// 	return
			// }
		}
		return &raft_cluster_v1.ElectionDecision{Term: r.Term}, nil
	} else {
		return &raft_cluster_v1.ElectionDecision{Term: r.Term}, errors.New("qourum not satisfied")
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
	r.ResetElectionTimer(ctx)

	commit := make(chan error, 1)
	go func() {

		r.mu.Lock()
		defer r.mu.Unlock()

		r.LeadId = int(req.IdLeader)
		r.Term = req.Term
		newLog := make([]model.Instance, 0)
		if len(req.Logs) != 0 {
			for _, log := range req.Logs {
				UUID, _ := uuid.Parse(log.Id)
				newLog = append(newLog, model.Instance{
					Id:      UUID,
					Content: model.JsonData{Name: log.JsonString},
					Term:    log.Term,
				})
			}
		}
		r.Logs = newLog
		if len(newLog) > 0 {
			Log.Debug("Updated", slog.Int("nodeID", int(r.IdNode)), slog.String("last", r.Logs[0].Content.Name))
		}
		r.SizeLogs = len(r.Logs)

		commit <- nil
	}()

	r.ResetElectionTimer(ctx)

	// control for context
	select {
	case <-ctx.Done():
		return &raft_cluster_v1.LeadAccept{Term: r.Term}, context.Canceled
	case err := <-commit:
		Log.Debug("Lead set", slog.String("port", r.Settings.Port), slog.Int("leadID", int(req.IdLeader)))
		return &raft_cluster_v1.LeadAccept{Term: r.Term}, err
	}
}

// Requesting vote from Candidate to Follower
func (r *ClusterNodeServer) RequestVote(ctx context.Context, req *raft_cluster_v1.RequestVoteRequest) (*raft_cluster_v1.RequestVoteResponse, error) {
	select {
	case <-ctx.Done():
		return &raft_cluster_v1.RequestVoteResponse{Term: r.Term}, context.Canceled
	default:
	}

	r.ResetElectionTimer(r.nodeCtx)
	defer r.ResetElectionTimer(r.nodeCtx)

	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.Term < req.Term {
		if r.State == Lead || r.State == Candidate {
			r.State = Follower
		}
		return &raft_cluster_v1.RequestVoteResponse{Term: r.Term}, nil
	} else if r.Term > req.Term {
		return &raft_cluster_v1.RequestVoteResponse{Term: r.Term}, errors.New("voter's term greater, candidate not legitimate")
	} else {
		// if request was sendet to high level nodes, they changing their state
		// it doesn`t mean that they always responsed 'yes', need check relevance logs
		if r.SizeLogs == 0 {
			return &raft_cluster_v1.RequestVoteResponse{Term: r.Term}, nil
		}

		// check relevance logs
		if req.LastLogTerm > r.Logs[r.SizeLogs-1].Term {
			// term more than our -> relevated, vote yes
			return &raft_cluster_v1.RequestVoteResponse{Term: r.Term}, nil
		} else if req.LastLogTerm < r.Logs[r.SizeLogs-1].Term {
			return &raft_cluster_v1.RequestVoteResponse{Term: r.Term}, errors.New("voter's term last log greater, candidate not legitimate")

		} else {
			// if terms equal checking lenght of logs
			if req.LastLogIndex > int64(r.SizeLogs)-1 {
				if r.State == Lead || r.State == Candidate {
					r.State = Follower
				}
				return &raft_cluster_v1.RequestVoteResponse{Term: r.Term}, nil
			} else {
				return &raft_cluster_v1.RequestVoteResponse{Term: r.Term}, errors.New("voter's logs more complete, candidate not legitimate")
			}
		}
	}
}

/*
----------------CRUD-------------------
*/

// Writing data to the node storage
func (r *ClusterNodeServer) LoadLog(ctx context.Context, req *raft_cluster_v1.LogInfo) (*raft_cluster_v1.LogAccept, error) {
	r.mu.RLock()
	// check node role
	if r.State == Candidate {
		return &raft_cluster_v1.LogAccept{Term: r.Term}, errors.New("forbidden, not saved")
	}
	if r.Term > req.Term {
		return &raft_cluster_v1.LogAccept{Term: r.Term}, errors.New("leader is not legitimate, not saved")
	}
	r.mu.RUnlock()

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

		r.mu.Lock()
		r.Logs = slices.Insert(r.Logs, int(req.Index), model.Instance{
			Id:      uuid.New(),
			Content: model.JsonData{Name: req.JsonString},
			Term:    req.Term,
		})
		r.mu.Unlock()
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

			Log.Info("Load log to Follower", slog.Int("nodeID", r.IdNode) /*slog.String("log", r.Logs[r.SizeLogs-1].Content.Name)*/)
			fmt.Println(r.Logs)
			r.Term = req.Term
			r.SizeLogs = len(r.Logs)
			return &raft_cluster_v1.LogAccept{Term: r.Term}, nil
		}
	}
}

// Load log to Cluster. This is an abstract method for changing data in the entire cluster.
func (r *ClusterNodeServer) Append(ctx context.Context, req *raft_cluster_v1.LogLeadRequest) (*raft_cluster_v1.Empty, error) {
	if r.State != Lead {
		Log.Info("Redirect to LEADNODE", slog.Int("LeadID", r.LeadId))
		_, err := r.Network[r.LeadId].Append(r.nodeCtx, req)
		return &raft_cluster_v1.Empty{}, err
	}
	if req == nil {
		return &raft_cluster_v1.Empty{}, errors.New("nil reference req")
	}
	log := &raft_cluster_v1.LogInfo{
		Id:         req.Id,
		Term:       r.Term,
		Index:      int64(r.SizeLogs),
		JsonString: req.JsonString,
	}
	_, err := r.LoadLog(r.nodeCtx, log)

	Log.Info("Load log to Lead", slog.Int("nodeID", r.IdNode), slog.String("log", r.Logs[r.SizeLogs-1].Content.Name))
	fmt.Println(r.Logs)
	if err != nil {
		return &raft_cluster_v1.Empty{}, err
	}

	// replication process
	loaded := 1
	for id, node := range r.Network {
		if loaded == r.Settings.Quorum {
			break
		}
		resp, err := node.LoadLog(r.nodeCtx, log)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				r.Network = slices.Delete(r.Network, id-1, id)
				continue
			} else {
				if resp.Term > r.Term {
					r.State = Follower
				}
				r.Logs = r.Logs[:len(r.Logs)-1]
				// TODO: rollback our log form prev nodes.
				return &raft_cluster_v1.Empty{}, err
			}
		}
		loaded++
	}
	if loaded < r.Settings.Quorum {
		// TODO: Delete log if not required
		return &raft_cluster_v1.Empty{}, errors.New("quorum not required")
	} else {
		return &raft_cluster_v1.Empty{}, nil
	}
}
