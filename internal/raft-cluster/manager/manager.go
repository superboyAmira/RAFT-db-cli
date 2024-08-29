package manager

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	cluster_node "warehouse/internal/raft-cluster/node"
)

const (
	cfgPath = "../../../configs/config.json"
)

type Manager struct {
	cluster     []*cluster_node.ClusterNodeServer
	activeNodes []*cluster_node.ClusterNodeServer
	leadId      int

	stop chan struct{}
	wg   sync.WaitGroup
	mu   sync.Mutex
}

func (r *Manager) StartCluster() error {
	if len(r.cluster) > 0 {
		return errors.New("alredy running")
	}

	settings, err := loadCfg()
	if err != nil {
		return err
	}

	// crete our cluster with network connection
	r.cluster = make([]*cluster_node.ClusterNodeServer, settings.Size)
	r.leadId = 0
	for i, state := 0, cluster_node.Lead; i < settings.Size; i++ {
		r.cluster[i] = cluster_node.New(i, state, r.leadId, nil, settings)
		state = cluster_node.Follower
	}
	for i, node := range r.cluster {
		for i_cl := 0; i < settings.Size; i_cl++ {
			if i_cl != i {
				node.Network = append(node.Network, r.cluster[i_cl])
			}
		}
	}

	// start all nodes
	r.activeNodes = make([]*cluster_node.ClusterNodeServer, 0)

	r.mu = sync.Mutex{}
	r.stop = make(chan struct{}, 1)
	for _, node := range r.cluster {
		r.wg.Add(1)
		go func() {
			defer r.wg.Done()
			err := node.Serve(r.stop)
			if err == nil {
				r.mu.Lock()
				r.activeNodes = append(r.activeNodes, node)
				r.mu.Unlock()
			}
		}()
	}
	// check minimal cluster size requierment
	if len(r.activeNodes) < settings.ReplicationFactor {
		r.stop <- struct{}{}
		return errors.New("interal server error with nodes, replication factor not required")
	}
	return nil
}

func (r *Manager) GracefullyStop() {
	r.stop <- struct{}{}
	r.wg.Wait()
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
