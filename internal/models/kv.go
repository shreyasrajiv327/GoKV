package models

type PutRequest struct{
	Key string  `json:"key"`
	Value string `json:"value"`
}

type GetResponse struct{
	Key string `json:"key"`
	Value string `json:"value"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}