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

	cluster     []*cluster_node.ClusterNodeServer
	leadId      int

	ctx context.Context
	stopCluster context.CancelFunc
	wg  sync.WaitGroup
	mu  sync.Mutex
}

func (r *Manager) StartCluster(log *slog.Logger) error {
	if len(r.cluster) > 0 {
		return errors.New("alredy running")
	}
	log.Debug("Starting cluster...")

	settings, err := loadCfg()
	if err != nil {
		return err
	}

	// crete our cluster
	r.cluster = make([]*cluster_node.ClusterNodeServer, settings.Size)
	r.leadId = 0

	r.mu = sync.Mutex{}
	r.ctx, r.stopCluster = context.WithCancel(context.Background())
	countActive := 0

	// we have 2 fields state and initTerm for lead primary settings, without primary election
	for i, state, initTerm := 0, cluster_node.Lead, 1; i < settings.Size; i++ {
		port, err := getFreePort()
		if err != nil {
			continue
		}
		tSettings := *settings
		tSettings.Port = ":" + strconv.FormatInt(int64(port), 10)
		r.cluster[i] = cluster_node.New(i, state, r.leadId, initTerm, &tSettings)
		r.wg.Add(1)
		go func(id int) {
			defer r.wg.Done()
			err := r.cluster[i].Serve(r.ctx, log)
			if err != nil {
				r.mu.Lock()
				countActive--
				r.mu.Unlock()
			}
		}(i)
		// change fields for followers
		state = cluster_node.Follower
		initTerm = 0
		countActive++
		// waiting for the node to start so that the system issues a new port
		time.Sleep(5*time.Millisecond)
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
					log.Warn("Failed to create gRPC client",  slog.Int("node_id", acive_node.IdNode), slog.String("error", err.Error()))
				}
                client := raft_cluster_v1.NewClusterNodeClient(conn)
                node.Network = append(node.Network, client)
            }
        }
    }

	// create client for operations with cluster per Leader
	conn, err := grpc.NewClient("localhost"+r.cluster[r.leadId].Settings.Port,
	grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Warn("Failed to create gRPC client",  slog.Int("node_id", r.leadId), slog.String("error", err.Error()))
		r.GracefullyStop()
	}
	r.globalClient = raft_cluster_v1.NewClusterNodeClient(conn)

	log.Debug("Cluster started")

	return nil
}

func (r *Manager) GracefullyStop() {
	r.stopCluster()
	r.wg.Wait()
}

func (r *Manager) SetLog(uuid, jsonString string, log *slog.Logger) (error) {
	if len(r.cluster) == 0 {
		return errors.New("didn't running")
	}
	req := &raft_cluster_v1.LogLeadRequest{
		Id: uuid,
		JsonString: jsonString,
	}

	_, err := r.globalClient.Append(r.ctx, req)
	return err
}

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