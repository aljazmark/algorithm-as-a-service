package models

//Algo respone object structure
type AlgoResponse struct {
	Algorithm     string   `json:"algorithm"`
	Input         string   `json:"input"`
	Parameters    []string `json:"parameters,omitempty"`
	Output        string   `json:"output"`
	ExecutionTime string   `json:"executiontime"`
}

//Algo request object structure
type AlgoRequest struct {
	Input      string   `json:"input"`
	Parameters []string `json:"parameters,omitempty"`
}
