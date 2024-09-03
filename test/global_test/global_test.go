package manager_test

import (
	"testing"
	"time"
	"warehouse/internal/raft-cluster/manager"

	"github.com/google/uuid"
)

func TestBasicStartCluster(t *testing.T) {
	man := manager.Manager{}
	err := man.StartCluster()
	if err != nil {
		t.Error(err.Error())
	}

	time.Sleep(20 *time.Millisecond)
	man.GracefullyStop()
}

func TestLoadLog(t *testing.T) {
	man := manager.Manager{}
	err := man.StartCluster()
	if err != nil {
		t.Error(err.Error())
	}
	id := uuid.NewString()
	err = man.SetLog(id, "{\"name\": \"Chapayev Mustache comb\"}")
	if err != nil {
		t.Error(err.Error())
	}
	logs := man.GetLogs(0)
	if logs == nil {
		t.Error("logs nil responsed")
	} else if logs[0].Id.String() == id {
		t.Errorf("%s not equal %s", logs[0].Id.String(), id)
	}
	time.Sleep(1 * time.Second)

	man.GracefullyStop()
}

func TestNetworkErrWith1Node(t *testing.T) {
	man := manager.Manager{}
	err := man.StartTestCluster()
	if err != nil {
		t.Error(err.Error())
	}
	time.Sleep(5 * time.Millisecond)
	// lead died
	man.TestStopNodeId0()
	time.Sleep(30 * time.Millisecond)

	man.GracefullyStop()
}

func TestNetworkErrWith1NodeWithReplica(t *testing.T) {
	man := manager.Manager{}
	err := man.StartTestCluster()
	if err != nil {
		t.Error(err.Error())
	}
	// load log
	id := uuid.NewString()
	err = man.SetLog(id, "{\"name\": \"Chapayev Mustache comb\"}")
	if err != nil {
		t.Error(err.Error())
	}

	time.Sleep(5 * time.Millisecond)
	// lead died
	man.TestStopNodeId0()
	time.Sleep(15 * time.Millisecond)

	// now we randomly calculate the timeout leader, so we need to look at the logs for a correct test

	// // check log 1 node
	// logs := man.GetLogs(1)
	// if logs == nil {
	// 	t.Error("logs nil responsed")
	// } else if logs[0].Id.String() == id {
	// 	t.Errorf("%s not equal %s", logs[0].Id.String(), id)
	// }
	// // must be replicated to the 2 node
	// logs = man.GetLogs(2)
	// if logs == nil {
	// 	t.Error("logs nil responsed")
	// } else if logs[0].Id.String() == id {
	// 	t.Errorf("%s not equal %s", logs[0].Id.String(), id)
	// }
	// t.Log(logs)

	man.GracefullyStop()
}


func TestNa(t *testing.T) {
	man := manager.Manager{}
	err := man.StartTestCluster()
	if err != nil {
		t.Error(err.Error())
	}
	// load log
	id := uuid.NewString()
	err = man.SetLog(id, "{\"name\": \"Chapayev Mustache comb\"}")
	if err != nil {
		t.Error(err.Error())
	}

	time.Sleep(5 * time.Millisecond)
	// lead died
	man.TestStopNodeId0()
	time.Sleep(50 * time.Millisecond)

	// man.TestStopNodeId1()
	time.Sleep(15 * time.Millisecond)

	man.GracefullyStop()
}