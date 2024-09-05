package repl

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"warehouse/internal/raft-cluster/manager"

	"github.com/google/uuid"
)

const (
	intput    = "> "
	connect   = "Connected to a database of Warehouse 13 at %s\n"
	reconnect = "Reconnected to a database of Warehouse 13 at %s"
	nodes     = "Known nodes:\n%s"
)

type REPL struct {
	Client *manager.Manager
}

func (r *REPL) Exec() {
	r.Client = &manager.Manager{}
	r.Client.StartCluster()

	// for testing some background errors
	go func() {
		time.Sleep(40 * time.Second)
		r.Client.StopConcreteNode(0)
	}()

	go func() {
		time.Sleep(60 * time.Second)
		r.Client.StopConcreteNode(2)
	}()

	scanner := bufio.NewScanner(os.Stdin)

	lead, err := r.Client.GetLeadAddr()
	if err != nil {
		log.Fatalf(err.Error())
	}
	conNodes, err := r.Client.GetNodes()
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf(connect, lead)
	fmt.Printf(nodes, conNodes)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 3)
		command := strings.ToUpper(parts[0])

		switch command {
		case "HELP":
			printHelp()
		case "EXIT":
			r.Client.GracefullyStop()
			fmt.Println("Cluster stopped.")
			fmt.Println("Exiting REPL...")
			return
		case "SET":
			if len(parts) != 3 {
				fmt.Println("Usage: SET <uuid> <json>")
				continue
			}
			_, err :=uuid.Parse(parts[1])
			if err != nil {
				fmt.Println("Error UUID4:", err)
				continue
			}
			if err := r.Client.SetLog(parts[1], parts[2]); err != nil {
				fmt.Println("Error setting log:", err)
			} else {
				fmt.Printf("Created (%v replicas)\n", r.Client.Params.ReplicationFactor)
			}
		case "GET":
			if len(parts) != 2 {
				fmt.Println("Usage: GET <uuid>")
				continue
			}
			result, err := r.Client.GetLog(parts[1])
			if err != nil {
				fmt.Println("Error getting log:", err)
			} else if result == "" {
				fmt.Println("Log entry not found.")
			} else {
				fmt.Println(result)
			}
		case "DELETE":
			if len(parts) != 2 {
				fmt.Println("Usage: DELETE <uuid>")
				continue
			}
			if err := r.Client.DeleteLog(parts[1]); err != nil {
				fmt.Println("Error deleting log:", err)
			} else {
				fmt.Println("Log entry deleted.")
			}
		default:
			fmt.Println("Unknown command:", command)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading input:", err)
	}
}

func printHelp() {
	fmt.Println(`
Available commands:
  HELP                - Show this help message
  EXIT                - Exit REPL
  SET <uuid> <json>   - Set a log entry
  GET <uuid>          - Get a log entry
  DELETE <uuid>       - Delete a log entry
`)
}
