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

	time.Sleep(10 * time.Millisecond)
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
	} else if logs[0].Id != id {
		t.Errorf("%s not equal %s", logs[0].Id, id)
	}
	time.Sleep(10 * time.Millisecond)

	man.GracefullyStop()
}

func TestGetLog(t *testing.T) {
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
	uuid, err := man.GetLog(id)
	if err != nil {
		t.Error(err.Error())
	} else if uuid != "{\"name\": \"Chapayev Mustache comb\"}" {
		t.Errorf("%s not equal %s", uuid, id)
	}
	time.Sleep(10 * time.Millisecond)

	man.GracefullyStop()
}

func TestDeleteLog(t *testing.T) {
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
	} else if logs[0].Id != id {
		t.Errorf("%s!%s!", logs[0].Id, id)
	}
	time.Sleep(10 * time.Millisecond)

	err = man.DeleteLog(logs[0].Id, "")
	if err != nil {
		t.Error(err.Error())
	}

	time.Sleep(10 * time.Millisecond)

	man.GracefullyStop()
}

func TestNetworkErrWith1Node(t *testing.T) {
	man := manager.Manager{}
	err := man.StartCluster()
	if err != nil {
		t.Error(err.Error())
	}
	time.Sleep(5 * time.Millisecond)
	// lead died

	man.StopCluster[0]()

	time.Sleep(100 * time.Millisecond)

	man.GracefullyStop()
}

func TestNetworkErrWithLeadNodeWithReplica(t *testing.T) {
	man := manager.Manager{}
	err := man.StartCluster()
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
	man.StopCluster[0]()
	time.Sleep(30 * time.Millisecond)

	// now we randomly calculate the timeout leader, so we need to look at the logs for a correct test

	man.GracefullyStop()
}

func TestNetworkErrWithNotLeadNodeWithReplica(t *testing.T) {
	man := manager.Manager{}
	err := man.StartCluster()
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
	man.StopCluster[1]()
	time.Sleep(30 * time.Millisecond)

	// now we randomly calculate the timeout leader, so we need to look at the logs for a correct test

	man.GracefullyStop()
}

func TestTwoDeactivated(t *testing.T) {
	man := manager.Manager{}
	err := man.StartCluster() // for connection global ctx
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
	man.StopCluster[0]()
	time.Sleep(50 * time.Millisecond)

	man.StopCluster[1]()
	time.Sleep(50 * time.Millisecond)

	man.GracefullyStop()
}
