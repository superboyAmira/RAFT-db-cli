package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"warehouse/internal/model"
	cluster_node "warehouse/internal/raft-cluster/node"
	"warehouse/pkg/raft/raft_cluster_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	cfgPath = "../../configs/config.json"
)

type Manager struct {
	Params       *cluster_node.ClusterSettings
	globalClient raft_cluster_v1.ClusterNodeClient
	globalId     int

	cluster []*cluster_node.ClusterNodeServer
	leadId  int

	ctx         []context.Context
	StopCluster []context.CancelFunc

	activeCluster bool
	wg            sync.WaitGroup
	mu            sync.Mutex
}

func (r *Manager) StartCluster() error {
	if len(r.cluster) > 0 {
		return errors.New("alredy running")
	}
	cluster_node.Log.Debug("Starting cluster...")

	var err error
	r.Params, err = loadCfg()
	if err != nil {
		return err
	}

	// crete our cluster
	r.cluster = make([]*cluster_node.ClusterNodeServer, r.Params.Size)
	r.ctx = make([]context.Context, r.Params.Size)
	r.StopCluster = make([]context.CancelFunc, r.Params.Size)
	r.leadId = -1

	r.mu = sync.Mutex{}

	for i := 0; i < r.Params.Size; i++ {
		r.ctx[i], r.StopCluster[i] = context.WithCancel(context.Background())
	}
	countActive := 0

	for i := 0; i < r.Params.Size; i++ {
		port, err := getFreePort()
		if err != nil {
			continue
		}
		tSettings := *r.Params
		tSettings.Port = ":" + strconv.FormatInt(int64(port), 10)
		r.cluster[i] = cluster_node.New(i, cluster_node.Follower, r.leadId, 0, &tSettings)
		r.wg.Add(1)
		go func(id int) {
			defer r.wg.Done()
			err := r.cluster[i].Serve(r.ctx[i])
			if err != nil {
				r.mu.Lock()
				countActive--
				r.mu.Unlock()
			}
		}(i)
		// change fields for followers
		countActive++
		// waiting for the node to start so that the system issues a new port
		time.Sleep(5 * time.Millisecond)
	}

	// check minimal cluster size requierment
	if countActive < r.Params.ReplicationFactor {
		r.GracefullyStop()
		return errors.New("interal server error with nodes, replication factor not required")
	}

	// create Network
	for i, node := range r.cluster {
		for i_cl, acive_node := range r.cluster {
			if i_cl != i {
				conn, err := grpc.NewClient(node.Settings.Host+acive_node.Settings.Port,
					grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					cluster_node.Log.Warn("Failed to create gRPC client", slog.Int("node_id", acive_node.IdNode), slog.String("error", err.Error()))
				}
				client := raft_cluster_v1.NewClusterNodeClient(conn)
				node.Network = append(node.Network, client)
			}
		}

	}

	// create client for operations with cluster node 0.
	// In CRUD operations we provided for redirect a request to current lead node
	r.globalId = 0
	conn, err := grpc.NewClient(r.cluster[0].Settings.Host+r.cluster[0].Settings.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cluster_node.Log.Warn("Failed to create gRPC client", slog.Int("node_id", r.leadId), slog.String("error", err.Error()))
		r.GracefullyStop()
	}
	r.globalClient = raft_cluster_v1.NewClusterNodeClient(conn)

	cluster_node.Log.Debug("Cluster started")

	r.activeCluster = true
	r.wg.Add(1)
	// to start searching for a working note for a client
	go r.pingActiveClient()

	time.Sleep(25 * time.Millisecond)
	return nil
}

// func to reconnected a client if client node isn't valid
func (r *Manager) pingActiveClient() {
	r.wg.Done()
	for {
		time.Sleep(3 * time.Second)
		if !r.activeCluster {
			return
		}

		// ping our lead
		_, err := r.globalClient.IsLead(r.ctx[r.globalId], nil)
		if err != nil {
			for i := 0; i < r.Params.Size; i++ {
				cluster_node.Log.Debug("try connect to alive node", slog.Int("i", i), slog.String("port", r.cluster[i].Settings.Port))
				if i > r.Params.ReplicationFactor-1 {
					fmt.Printf("cluster size (%v) is smaller than a replication factor (%v)!\n", r.Params.Size-i, r.Params.ReplicationFactor)
					r.GracefullyStop()
					return
				}
				conn, err := grpc.NewClient(r.cluster[i].Settings.Host+r.cluster[i].Settings.Port,
					grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					cluster_node.Log.Warn("Failed to create gRPC client", slog.Int("node_id", r.leadId), slog.String("error", err.Error()))
					continue
				}
				r.globalClient = raft_cluster_v1.NewClusterNodeClient(conn)
				_, err = r.globalClient.IsLead(r.ctx[i], nil)
				if err != nil {
					cluster_node.Log.Debug("client err", slog.String("err", err.Error()))
					continue
				}
				fmt.Printf("Reconnected to a database of Warehouse 13 at %s%s\n>", r.cluster[i].Settings.Host, r.cluster[i].Settings.Port)
				r.globalId = i
				break
			}
		}
	}
}

func (r *Manager) GracefullyStop() {
	r.activeCluster = false
	for _, cancel := range r.StopCluster {
		cancel()
	}
	r.wg.Wait()
}

func (r *Manager) StopConcreteNode(id int) {
	if id > len(r.cluster) || id < 0 {
		return
	}
	r.StopCluster[id]()
}

/*
----------------CRUD-------------------
*/

func (r *Manager) SetLog(uuid, jsonString string) error {
	if len(r.cluster) == 0 {
		return errors.New("didn't running")
	}
	req := &raft_cluster_v1.LogLeadRequest{
		Id:         uuid,
		JsonString: jsonString,
	}

	_, err := r.globalClient.Append(r.ctx[r.globalId], req)
	return err
}

func (r *Manager) DeleteLog(uuid string) error {
	if len(r.cluster) == 0 {
		return errors.New("didn't running")
	}
	req := &raft_cluster_v1.LogLeadRequest{
		Id: uuid,
	}

	_, err := r.globalClient.Delete(r.ctx[r.globalId], req)
	return err
}

func (r *Manager) GetLog(uuid string) (string, error) {
	if len(r.cluster) == 0 {
		return "", errors.New("didn't running")
	}
	req := &raft_cluster_v1.LogLeadRequest{
		Id:         uuid,
		JsonString: "",
	}

	resp, err := r.globalClient.Get(r.ctx[r.globalId], req)
	if err != nil {
		return "", err
	}
	return resp.JsonString, nil
}

func (r *Manager) GetLogs(nodeID int) []model.Instance {
	return r.cluster[nodeID].Logs
}

/*
----------------Support-------------------
*/

func loadCfg() (*cluster_node.ClusterSettings, error) {
	file, err := os.Open(cfgPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &cluster_node.ClusterSettings{}

	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	config.HeartBeatIntervalLeader = 1 * time.Millisecond
	config.HeartBeatTimeout = 3 * time.Millisecond
	config.ElectionTimeout = 2 * time.Millisecond
	return config, nil
}

func getFreePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

/*
----------------Client REPL Info-------------------
*/

func (r *Manager) GetLeadAddr() (string, error) {
	if len(r.cluster) == 0 {
		return "", errors.New("didn't running")
	}
	for {
		for i := 0; i < len(r.cluster); i++ {
			if r.cluster[i].State == cluster_node.Lead {
				conn, err := grpc.NewClient(r.cluster[i].Settings.Host+r.cluster[i].Settings.Port,
					grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					cluster_node.Log.Warn("Failed to create gRPC client", slog.Int("node_id", r.leadId), slog.String("error", err.Error()))
					r.GracefullyStop()
					return "", err
				}
				r.globalClient = raft_cluster_v1.NewClusterNodeClient(conn)
				return "localhost" + r.cluster[i].Settings.Port, nil
			}
		}
	}
}

func (r *Manager) GetNodes() (string, error) {
	if len(r.cluster) == 0 {
		return "", errors.New("didn't running")
	}
	res := strings.Builder{}
	for _, node := range r.cluster {
		res.WriteString(node.Settings.Host + node.Settings.Port + "\n")
	}
	return res.String(), nil
}
