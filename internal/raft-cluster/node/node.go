package node

import (
	"fmt"
	"warehouse/internal/model"
)

const (
	Follower = iota
	Candidate
	Lead
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
type ClusterNode struct {
	idNode int	
	status   bool
	settings *config

	log []model.Instance

	idLead int
}

func NewClusterNode(id int) *ClusterNode {

	return &ClusterNode{
		idNode: id,
		status: true,
		settings: load(),

		log: make([]model.Instance, 0),

	}
}

func (r *ClusterNode) run() {}

// send HeartBeat string to global monitoring system
func (r *ClusterNode) GetHeartBeat() string {
	res := fmt.Sprintf("%d", r.settings.Replication) + r.settings.Port
	if !r.status {
		res += " stopped"
	}
	return res
}
