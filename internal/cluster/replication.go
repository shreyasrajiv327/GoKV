package cluster

import(
	"bytes"
	"encoding/json"
	"net/http"
)

type Replicator struct{
	peers []string
}

func NewReplicator(peers []string) *Replicator{
	return &Replicator{
		peers: peers,
	}
}

type ReplicationRequest struct{
	Key string `json:"key"`
	Value string `json:"value"`
}

func(r *Replicator) ReplicatePut(key, value string){
	body := ReplicationRequest{
		Key: key,
		Value: value,
	}

	jsonData, _ := json.Marshal(body)

	for _, peer := range r.peers{

		http.Post(
			peer+"/replicate/put",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
	}
}

func (r *Replicator) ReplicateDelete(key string) {

	client := &http.Client{}

	for _, peer := range r.peers {

		req, _ := http.NewRequest(
			http.MethodDelete,
			peer+"/replicate/delete/"+key,
			nil,
		)

		client.Do(req)
	}
}