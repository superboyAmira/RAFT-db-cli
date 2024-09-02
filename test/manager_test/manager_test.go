package manager_test

import (
	"testing"
	"time"
	"warehouse/internal/raft-cluster/manager"
)

func TestBasicStartCluster(t *testing.T) {
	man := manager.Manager{}
	// log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	err := man.StartCluster()
	if err != nil {
		t.Error(err.Error())
	}
	man.GracefullyStop()
}

func TestLoadLog(t *testing.T) {
	man := manager.Manager{}
	// log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	err := man.StartCluster()
	if err != nil {
		t.Error(err.Error())
	}
	time.Sleep(29 * time.Second)
	// err = man.SetLog(uuid.NewString(), "{\"name\": \"Chapayev Mustache comb\"}")
	// if err != nil {
	// 	t.Error(err.Error())
	// }

	man.GracefullyStop()
}

///time=2024-09-01T23:07:53.351+03:00 level=DEBUG msg="Election result" port=:46017 err="rpc error: code = Unknown desc = lead isn't legitimate"