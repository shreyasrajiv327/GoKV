package models

type JoinRequest struct {
    NodeID   string `json:"node_id"`
    RaftAddr string `json:"raft_addr"`
}

type RemoveRequest struct {
    NodeID string `json:"node_id"`
}