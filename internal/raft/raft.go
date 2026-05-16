package raft

import(
	"fmt"
	"gokv/internal/cache"
	"os"
	"path/filepath"
	"time"

	hashiraft "github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
)

func SetupRaft(
	nodeID string,
	raftDir string,
	raftBind string,
	store *cache.MemoryStore,
)(*hashiraft.Raft, error) {

	config := hashiraft.DefaultConfig()
	config.LocalID = hashiraft.ServerID(nodeID)



transport, err := hashiraft.NewTCPTransport(
	raftBind,
	nil,
	3,
	10*time.Second,
	os.Stdout,
)

if err != nil {
	return nil, err
}
	snapshotStore, err := hashiraft.NewFileSnapshotStore(
		raftDir,
		2,
		os.Stdout,
	)

	if err != nil {
		return nil, err
	}

	logStore, err := boltdb.NewBoltStore(
		filepath.Join(raftDir, "raft-log.bolt"),
	)

	if err != nil {
		return nil, err
	}

	stableStore, err := boltdb.NewBoltStore(
		filepath.Join(raftDir, "raft-stable.bolt"),
	)

	if err != nil {
		return nil, err
	}

	fsm := NewFSM(store)

	raftNode, err := hashiraft.NewRaft(
		config,
		fsm,
		logStore,
		stableStore,
		snapshotStore,
		transport,
	)

	if err != nil {
		return nil, err
	}

	bootstrapConfig := hashiraft.Configuration{
		Servers: []hashiraft.Server{
			{
				ID:      config.LocalID,
				Address: transport.LocalAddr(),
			},
		},
	}

	raftNode.BootstrapCluster(bootstrapConfig)

	fmt.Println("Raft node initialized:", nodeID)

	return raftNode, nil
	
}