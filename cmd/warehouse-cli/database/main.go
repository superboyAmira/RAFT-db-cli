package main

import (
	"sync"
	"warehouse/internal/raft-cluster/manager"
)


func main() {
	man := manager.Manager{}
	// log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := man.StartCluster()
		if err != nil {
			// t.Error(err.Error())
			return
		}
	}()
	wg.Wait()
}