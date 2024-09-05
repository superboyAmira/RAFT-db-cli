# RAFT-db-cli
Fault-tolerant distributed local database system using the RAFT algorithm. Golang

## Tech Stack

* Go 1.22.3
* GRPC + protobuf

# Resources (links)
 
* https://www.youtube.com/watch?v=64Zp3tzNbpE MIT course
* go implemenation of RAFT protocol https://github.com/hashicorp/raft
* https://medium.com/@govinda.attal/raft-consensus-leader-election-with-golang-89bfdbd471cb

# About project

## RAFT (Replicated and Fault-Tolerant)
RAFT is a consensus algorithm designed to simplify the implementation of distributed systems that require agreement, such as distributed databases or storage systems.
RAFT was proposed in 2014 by Diego Ongaro and John Ousterhout as an alternative to the more complex Paxos algorithm.

### Key Principles of the RAFT Protocol:
**Leader Election:**

* In RAFT, there is always one node acting as the leader, while the other nodes are followers. The leader is responsible for managing the system and replicating logs.
* If the leader does not respond within a specified timeout, a new election process begins, and one of the followers becomes the leader.

**Log Replication:**

* The leader accepts requests from clients and records them in its log, then replicates these entries to the followers' logs.
* Once an entry is replicated on a majority of nodes (quorum), it is considered committed, and the leader can inform the client of the operation's completion.

**Safety:**

* RAFT ensures that once a node confirms an operation, it will not be lost, even if the leader fails.
* A new leader is chosen in such a way that all previously confirmed entries remain consistent.

**Consensus:**

* All nodes in the system must agree on the order of command execution. This is achieved through leader election and log replication mechanisms.

### Rules for Ensuring Fault Tolerance:
* **Quorum:**

  * A decision requires the approval of a majority of nodes, meaning the system can continue to operate even if some nodes fail.
  
* **State Persistence:**

  * Each node saves its state to disk, so in case of a failure, it can recover and continue from the same point.

* **Failure Isolation:**

  * If one node fails, the others continue to function until the quorum is restored.

* **Failure Detection and Recovery:**

  * The system automatically detects failures and initiates a new leader election or data recovery.

# Quick start

To verify compliance with the RAFT principle. the global cluster logger should be changed to Info Level. and see the result in the logs.This is a temporary solution, in the future there will be a more visual and easy-to-analyze test.

In the current implementation, all possible cases of election and recovery of node network errors have been handled. There is functionality for adding, deleting and viewing database records.

For more detail information of protocol, you can set Logger level to Debug. (you need to perform it in a short time interval and write it to a file, since there will be a lot of logs)

to start, use
```cd cmd/warehouse-cli/```
```go run main.go```

you can write 'help or HELP' in REPL to see the possible commands

* the current tests are located in the global test folder

# Demonstration of the functionality

Example of Set and Get Operations
(with Info level)

```
➜  warehouse-cli git:(main) ✗ go run main.go
time=2024-09-03T09:31:58.370+03:00 level=INFO msg="Starting cluster..."
time=2024-09-03T09:31:58.370+03:00 level=INFO msg="Node is running on port" node_id=0 address=[::]:33267
time=2024-09-03T09:31:58.376+03:00 level=INFO msg="Node is running on port" node_id=1 address=[::]:44115
time=2024-09-03T09:31:58.381+03:00 level=INFO msg="Node is running on port" node_id=2 address=[::]:36535
time=2024-09-03T09:31:58.386+03:00 level=INFO msg="Cluster started"
Connected to a database of Warehouse 13 at localhost:33267
Known nodes:
localhost:33267
localhost:44115
localhost:36535
> SET 0d5d3807-5fbf-4228-a657-5a091c4e497f '{"name": "Chapayev's Mustache comb"}'
time=2024-09-03T09:32:23.352+03:00 level=INFO msg="Load log to Follower" nodeID=0 log=0d5d3807-5fbf-4228-a657-5a091c4e497f
time=2024-09-03T09:32:23.352+03:00 level=INFO msg="Load log to Lead" nodeID=0 log="'{\"name\": \"Chapayev's Mustache comb\"}'" log=0d5d3807-5fbf-4228-a657-5a091c4e497f
time=2024-09-03T09:32:23.352+03:00 level=INFO msg="Load log to Follower" nodeID=1 log=0d5d3807-5fbf-4228-a657-5a091c4e497f
Log entry set.
> GET 0d5d3807-5fbf-4228-a657-5a091c4e497f
Log entry: '{"name": "Chapayev's Mustache comb"}'
> EXIT
time=2024-09-03T09:32:54.241+03:00 level=INFO msg="Death note about node" node_id=2 leadID=0 term=1 address=[::]:36535
time=2024-09-03T09:32:54.241+03:00 level=INFO msg="Death note about node" node_id=1 leadID=0 term=1 last_content="'{\"name\": \"Chapayev's Mustache comb\"}'"
time=2024-09-03T09:32:54.241+03:00 level=INFO msg=Shutdown... node_id=1 address=[::]:44115
time=2024-09-03T09:32:54.241+03:00 level=INFO msg=Shutdown... node_id=2 address=[::]:36535
time=2024-09-03T09:32:54.241+03:00 level=INFO msg="Death note about node" node_id=0 leadID=0 term=1 last_content="'{\"name\": \"Chapayev's Mustache comb\"}'"
time=2024-09-03T09:32:54.242+03:00 level=INFO msg=Shutdown... node_id=0 address=[::]:33267
time=2024-09-03T09:32:54.242+03:00 level=INFO msg="Stopped Gracefully..." node_id=1 address=[::]:44115
time=2024-09-03T09:32:54.242+03:00 level=INFO msg="Stopped Gracefully..." node_id=2 address=[::]:36535
time=2024-09-03T09:32:54.243+03:00 level=INFO msg="Stopped Gracefully..." node_id=0 address=[::]:33267
Cluster stopped.
Exiting REPL...
```

Example of Set, Get, Delete Operations
(with Error(Prod) level)

```
➜  warehouse-cli git:(main) ✗ go run main.go
Connected to a database of Warehouse 13 at localhost:38503
Known nodes:
localhost:38503
localhost:40727
localhost:37167
> SET 429420942 incorrect
Error UUID4: invalid UUID length: 9
> SET 0d5d3807-5fbf-4228-a657-5a091c4e497f some test
Created (2 replicas)
> GET 0d5d3807-5fbf-4228-a657-5a091c4e497f
some test
> 
> DELETE 0d5d3807-5fbf-4228-a657-5a091c4e497f
Log entry deleted.
> SET 0d5d3807-5fbf-4228-a657-5a091c4e497f frefe
Created (2 replicas)
> 
> Reconnected to a database of Warehouse 13 at localhost:40727
>GET 0d5d3807-5fbf-4228-a657-5a091c4e497f
frefe
> cluster size (1) is smaller than a replication factor (2)!
...
```

**There is an understanding of how to rewrite the implementation in a more optimized way**

# thanks!