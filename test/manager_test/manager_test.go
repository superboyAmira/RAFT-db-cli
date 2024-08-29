package manager_test

import (
	"log/slog"
	"os"
	"testing"
	"warehouse/internal/raft-cluster/manager"
)


func TestBasicStartCluster(t *testing.T) {
	man := manager.Manager{}
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	err := man.StartCluster(log)
	if err != nil {
		t.Error(err.Error())
	}
	man.GracefullyStop()
}