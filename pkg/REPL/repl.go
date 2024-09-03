package repl

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"warehouse/internal/raft-cluster/manager"
)

const (
	intput = "> "
	connect = "Connected to a database of Warehouse 13 at %s\n"
	reconnect = "Reconnected to a database of Warehouse 13 at %s"
	nodes = "Known nodes:\n%s"
)

type REPL struct {
	Client *manager.Manager
	terminate chan os.Signal 
}

func (r *REPL) Exec() {
	r.Client = &manager.Manager{}
	r.Client.StartCluster()
	time.Sleep(3*time.Second)
	scanner := bufio.NewScanner(os.Stdin)

	r.terminate = make(chan os.Signal, 1)
	signal.Notify(r.terminate, os.Interrupt, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT, syscall.SIGTSTP)
	go func ()  {
		<-r.terminate
		r.Client.GracefullyStop()
	}()

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
			if err := r.Client.SetLog(parts[1], parts[2]); err != nil {
				fmt.Println("Error setting log:", err)
			} else {
				fmt.Println("Log entry set.")
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
				fmt.Println("Log entry:", result)
			}
		case "DELETE":
			if len(parts) != 3 {
				fmt.Println("Usage: DELETE <uuid> <json>")
				continue
			}
			if err := r.Client.DeleteLog(parts[1], parts[2]); err != nil {
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
  DELETE <uuid> <json> - Delete a log entry
`)
}