package node

type config struct {
	Replication int    `cfg:"replica"`
	Host        string `cfg:"host"`
	Port        string
}

func load() *config {
	// todo
	// random free port gen
	return nil
}
