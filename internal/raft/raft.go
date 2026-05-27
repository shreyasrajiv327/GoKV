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
) (*hashiraft.Raft,*FSM,error){

	config:=hashiraft.DefaultConfig()

	config.LocalID=
		hashiraft.ServerID(nodeID)

	transport,err:=
		hashiraft.NewTCPTransport(
			raftBind,
			nil,
			3,
			10*time.Second,
			os.Stdout,
		)

	if err!=nil{
		return nil,nil,err
	}

	snapshotStore,err:=
		hashiraft.NewFileSnapshotStore(
			raftDir,
			2,
			os.Stdout,
		)

	if err!=nil{
		return nil,nil,err
	}

	logStore,err:=
		boltdb.NewBoltStore(
			filepath.Join(
				raftDir,
				"raft-log.bolt",
			),
		)

	if err!=nil{
		return nil,nil,err
	}

	stableStore,err:=
		boltdb.NewBoltStore(
			filepath.Join(
				raftDir,
				"raft-stable.bolt",
			),
		)

	if err!=nil{
		return nil,nil,err
	}

	fsm:=NewFSM(store)

	raftNode,err:=
		hashiraft.NewRaft(
			config,
			fsm,
			logStore,
			stableStore,
			snapshotStore,
			transport,
		)

	if err!=nil{
		return nil,nil,err
	}

	existing, err :=
	hashiraft.HasExistingState(
		logStore,
		stableStore,
		snapshotStore,
	)

if err != nil {

	return nil,
		nil,
		err
}

if !existing && nodeID == "node1" {

	bootstrap :=
		hashiraft.Configuration{

			Servers:
				[]hashiraft.Server{

					{
						ID:
						hashiraft.ServerID(
							nodeID,
						),

						Address:
						hashiraft.ServerAddress(
							raftBind,
						),
					},
				},
		}

	future :=
		raftNode.
			BootstrapCluster(
				bootstrap,
			)

	if err :=
		future.Error(); err != nil {

		if err !=
			hashiraft.ErrCantBootstrap {

			return nil,
				nil,
				err
		}
	}
}

	if os.Getenv("BOOTSTRAP")=="true"{

		cfg:=hashiraft.Configuration{
			Servers:[]hashiraft.Server{
				{
					ID:hashiraft.ServerID(nodeID),
					Address:hashiraft.ServerAddress(
						raftBind,
					),
				},
			},
		}

		future:=
			raftNode.
			BootstrapCluster(
				cfg,
			)

		if err:=future.Error();err!=nil{

			if err!=
				hashiraft.ErrCantBootstrap{

				return nil,nil,err
			}
		}
	}

	fmt.Println(
		"raft initialized",
		nodeID,
	)

	return raftNode,fsm,nil
}