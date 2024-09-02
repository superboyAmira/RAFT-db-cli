package manager

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"os"
	"strconv"
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
	globalClient raft_cluster_v1.ClusterNodeClient

	cluster []*cluster_node.ClusterNodeServer
	leadId  int

	ctx         context.Context
	stopCluster context.CancelFunc
	wg          sync.WaitGroup
	mu          sync.Mutex

	// testing
	contexts []context.Context
	cancels  []context.CancelFunc
}

func (r *Manager) StartCluster() error {
	if len(r.cluster) > 0 {
		return errors.New("alredy running")
	}
	cluster_node.Log.Info("Starting cluster...")

	settings, err := loadCfg()
	if err != nil {
		return err
	}

	// crete our cluster
	r.cluster = make([]*cluster_node.ClusterNodeServer, settings.Size)
	r.leadId = -1

	r.mu = sync.Mutex{}
	r.ctx, r.stopCluster = context.WithCancel(context.Background())
	countActive := 0

	for i := 0; i < settings.Size; i++ {
		port, err := getFreePort()
		if err != nil {
			continue
		}
		tSettings := *settings
		tSettings.Port = ":" + strconv.FormatInt(int64(port), 10)
		r.cluster[i] = cluster_node.New(i, cluster_node.Follower, r.leadId, 0, &tSettings)
		r.wg.Add(1)
		go func(id int) {
			defer r.wg.Done()
			err := r.cluster[i].Serve(r.ctx)
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
	if countActive < settings.ReplicationFactor {
		r.GracefullyStop()
		return errors.New("interal server error with nodes, replication factor not required")
	}

	// create Network
	for i, node := range r.cluster {
		for i_cl, acive_node := range r.cluster {
			if i_cl != i {
				conn, err := grpc.NewClient("localhost"+acive_node.Settings.Port,
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
	conn, err := grpc.NewClient("localhost"+r.cluster[0].Settings.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cluster_node.Log.Warn("Failed to create gRPC client", slog.Int("node_id", r.leadId), slog.String("error", err.Error()))
		r.GracefullyStop()
	}
	r.globalClient = raft_cluster_v1.NewClusterNodeClient(conn)

	cluster_node.Log.Info("Cluster started")

	time.Sleep(5 * time.Millisecond)
	return nil
}

func (r *Manager) GracefullyStop() {
	r.stopCluster()
	r.wg.Wait()
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

	_, err := r.globalClient.Append(r.ctx, req)
	return err
}

func (r *Manager) GetLogs() []model.Instance {
	return r.cluster[0].Logs
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
----------------Testing-------------------
*/

func (r *Manager) StartTestCluster() error {
	if len(r.cluster) > 0 {
		return errors.New("alredy running")
	}
	cluster_node.Log.Info("TEST Starting cluster...")

	settings, err := loadCfg()
	if err != nil {
		return err
	}

	// init contexts

	for i := 0; i < settings.Size; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		r.contexts = append(r.contexts, ctx)
		r.cancels = append(r.cancels, cancel)
	}

	// crete our cluster
	r.cluster = make([]*cluster_node.ClusterNodeServer, settings.Size)
	r.leadId = -1

	r.mu = sync.Mutex{}
	r.ctx, r.stopCluster = context.WithCancel(context.Background())
	countActive := 0

	for i := 0; i < settings.Size; i++ {
		port, err := getFreePort()
		if err != nil {
			continue
		}
		tSettings := *settings
		tSettings.Port = ":" + strconv.FormatInt(int64(port), 10)
		r.cluster[i] = cluster_node.New(i, cluster_node.Follower, r.leadId, 0, &tSettings)
		r.wg.Add(1)
		go func(id int) {
			defer r.wg.Done()
			err := r.cluster[i].Serve(r.contexts[i])
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
	if countActive < settings.ReplicationFactor {
		r.GracefullyStop()
		return errors.New("interal server error with nodes, replication factor not required")
	}

	// create Network
	for i, node := range r.cluster {
		for i_cl, acive_node := range r.cluster {
			if i_cl != i {
				conn, err := grpc.NewClient("localhost"+acive_node.Settings.Port,
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
	conn, err := grpc.NewClient("localhost"+r.cluster[0].Settings.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cluster_node.Log.Warn("Failed to create gRPC client", slog.Int("node_id", r.leadId), slog.String("error", err.Error()))
		r.GracefullyStop()
	}
	r.globalClient = raft_cluster_v1.NewClusterNodeClient(conn)

	cluster_node.Log.Info("TEST Cluster started")

	time.Sleep(5 * time.Millisecond)
	return nil
}

func (r *Manager) TestStopNode(index int) {
	r.cancels[index]()
	cluster_node.Log.Debug("Stopped ", slog.Int("id", index))
}
