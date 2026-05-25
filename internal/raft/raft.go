package raft

import (
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
) (*hashiraft.Raft, *FSM, error) {

	config := hashiraft.DefaultConfig()

	config.LocalID =
		hashiraft.ServerID(nodeID)

	transport, err :=
		hashiraft.NewTCPTransport(
			raftBind,
			nil,
			3,
			10*time.Second,
			os.Stdout,
		)

	if err != nil {

		return nil,
			nil,
			err
	}

	snapshotStore, err :=
		hashiraft.NewFileSnapshotStore(
			raftDir,
			2,
			os.Stdout,
		)

	if err != nil {

		return nil,
			nil,
			err
	}

	logStore, err :=
		boltdb.NewBoltStore(
			filepath.Join(
				raftDir,
				"raft-log.bolt",
			),
		)

	if err != nil {

		return nil,
			nil,
			err
	}

	stableStore, err :=
		boltdb.NewBoltStore(
			filepath.Join(
				raftDir,
				"raft-stable.bolt",
			),
		)

	if err != nil {

		return nil,
			nil,
			err
	}

	fsm :=
		NewFSM(store)

	raftNode, err :=
		hashiraft.NewRaft(
			config,
			fsm,
			logStore,
			stableStore,
			snapshotStore,
			transport,
		)

	if err != nil {

		return nil,
			nil,
			err
	}

	bootstrapConfig :=
		hashiraft.Configuration{

			Servers:
			[]hashiraft.Server{

				{
					ID:
					"node1",

					Address:
					"127.0.0.1:7001",
				},

				{
					ID:
					"node2",

					Address:
					"127.0.0.1:7002",
				},

				{
					ID:
					"node3",

					Address:
					"127.0.0.1:7003",
				},
			},
		}

	if nodeID == "node1" {

		future :=
			raftNode.
			BootstrapCluster(
				bootstrapConfig,
			)

		err :=
			future.Error()

		if err != nil {

			if err !=
				hashiraft.
				ErrCantBootstrap {

				return nil,
					nil,
					err
			}
		}
	}

	fmt.Println(
		"Raft node initialized:",
		nodeID,
	)

	return raftNode,
		fsm,
		nil
}