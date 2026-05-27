# GoKV

GoKV is a distributed fault-tolerant in-memory key-value store built in Go using the Raft consensus algorithm.

It supports leader election, automatic failover, replicated state machines, dynamic cluster membership, persistence, snapshots, and multiple consistency guarantees.

The system is designed to explore real distributed systems concepts while remaining usable as a lightweight local distributed cache / key-value store.

---

## Features

### Distributed Consensus

- Raft consensus using HashiCorp Raft
- Automatic leader election
- Leader failover
- Replicated state machine (FSM)
- Quorum-based safety guarantees
- Dynamic cluster membership
- N-node cluster support
- Leader-only writes
- Automatic node recovery after restart

---

### Persistence

- Write Ahead Log (WAL)
- Snapshot persistence
- Raft persistent logs
- Automatic recovery after crash
- Replay WAL during startup
- Restore FSM state automatically

---

### Read Consistency Models

GoKV supports multiple consistency guarantees.

### Eventual Consistency

Reads local state immediately.

```http
GET /kv/user
```

Lowest latency.

---

### Strong Consistency

Read directly from leader.

```http
GET /kv/user?consistency=strong
```

Always returns latest committed value.

---

### Smart Consistency

Follower serves local read only when:

```text
applied_index >= leader_commit_index
```

Otherwise request automatically forwards to leader.

```http
GET /kv/user?consistency=smart
```

Balances correctness and latency.

---

## Dynamic Cluster Membership

Add nodes dynamically:

```http
POST /cluster/join
```

Remove nodes dynamically:

```http
DELETE /cluster/remove/:nodeID
```

Cluster membership changes are replicated through Raft.

Supports:

- Runtime scaling
- Quorum recalculation
- Safe node removal
- Leader failover after topology changes

---

## Architecture

```

                 Client
                    |
        +-----------+-----------+
        |                       |
        v                       v

      Leader Node         Follower Nodes
          |
          |
      Raft Replication
          |
          v

FSM (State Machine)
          |
          v

In Memory KV Store
          |
   +------+------+
   |             |

WAL Snapshot

```

---

## Tech Stack

- Go
- Gin
- HashiCorp Raft
- BoltDB
- WAL Persistence
- Snapshot Recovery

---

## Project Structure

```

GoKV/

├── cmd/server/
├── internal/
│
├── raft/
├── handlers/
├── services/
├── repository/
├── snapshot/
├── wal/
├── cache/
├── cluster/
│
├── raft-data-\*/
├── go.mod
└── README.md

```

---

## Running Locally

### 1. Start Bootstrap Node

Bootstrap node creates the initial cluster.

```bash
NODE_ID=node1 \
RAFT_PORT=7001 \
PORT=8001 \
go run cmd/server/main.go
```

Verify:

```bash
curl localhost:8001/node/status
```

Expected:

```json
{
"id":"node1",
"role":"leader"
}
```

---

### 2. Start Additional Nodes

Node 2:

```bash
NODE_ID=node2 \
RAFT_PORT=7002 \
PORT=8002 \
go run cmd/server/main.go
```

Join:

```bash
curl \
-X POST localhost:8001/cluster/join \
-H "Content-Type: application/json" \
-d '{
"node_id":"node2",
"raft_addr":"127.0.0.1:7002"
}'
```

---

Node 3:

```bash
NODE_ID=node3 \
RAFT_PORT=7003 \
PORT=8003 \
go run cmd/server/main.go
```

Join:

```bash
curl \
-X POST localhost:8001/cluster/join \
-H "Content-Type: application/json" \
-d '{
"node_id":"node3",
"raft_addr":"127.0.0.1:7003"
}'
```

---

## Verify Cluster

```bash
curl localhost:8001/node/status

curl localhost:8002/node/status

curl localhost:8003/node/status
```

Expected:

```text
1 leader
N-1 followers
Same term across cluster
```

---

## Example Usage

Write:

```bash
curl \
-X PUT localhost:8001/kv \
-H "Content-Type: application/json" \
-d '{"key":"user","value":"Tom"}'
```

Read:

```bash
curl localhost:8002/kv/user
```

Strong read:

```bash
curl "localhost:8002/kv/user?consistency=strong"
```

Smart read:

```bash
curl "localhost:8002/kv/user?consistency=smart"
```

Delete:

```bash
curl \
-X DELETE localhost:8001/kv/user
```

---

## Fault Tolerance Validation

GoKV supports:

### Leader Failure

```text
Leader crash
↓
Election timeout
↓
Majority vote
↓
New leader elected
```

---

### Node Restart Recovery

```text
Node crashes
↓
Restart process
↓
Load persisted raft state
↓
Rejoin cluster automatically
```

---

### Dynamic Membership Recovery

```text
Remove node
↓
Quorum recalculated
↓
Replication continues
```

---

## Cluster Safety Guarantees

Raft guarantees:

- No split brain
- Majority quorum required
- Leader-only writes
- Log consistency
- Safe failover
- Membership changes replicated through consensus

---

## Failure Tolerance

Raft cluster fault tolerance:

```text
1 node → tolerate 0 failures

3 nodes → tolerate 1 failure

5 nodes → tolerate 2 failures

7 nodes → tolerate 3 failures
```

Formula:

```text
Failures tolerated = floor((N - 1) / 2)
```

---

## Why FSM?

Raft replicates logs.

FSM applies committed logs into actual state.

Example:

```text
PUT user=shreyas

↓

Replicated Raft Log

↓

FSM Apply()

↓

Memory Store Updated
```

Without FSM:

```text
Consensus exists

But data never changes
```

FSM converts consensus into application state.

---

## Future Improvements

- Docker Compose deployment
- Metrics endpoint
- Graceful shutdown
- Authentication
- TLS
- Benchmark suite

---

## Learning Goals

GoKV explores:

- Consensus algorithms
- Distributed consistency
- Fault tolerance
- Quorum systems
- Replicated state machines
- Persistence and recovery
- Dynamic cluster membership
- Leader election

---

Built with Go and distributed systems principles.
