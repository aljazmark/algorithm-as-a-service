package models

//Help struct definition
type Help struct {
	Algorithms []string   `json:"algorithms"`
	Help       []AlgoHelp `json:"help"`
}

//AlgoHelp defines algorithm details response
type AlgoHelp struct {
	Algorithm         string `json:"algorithm"`
	Category          string `json:"category"`
	Description       string `json:"description"`
	InputFormat       string `json:"inputFormat"`
	InputExample      string `json:"inputExample"`
	ParametersFormat  string `json:"parametersFormat"`
	ParametersExample string `json:"parametersExample"`
	OutputFormat      string `json:"outputFormat"`
	OutputExample     string `json:"outputExample"`
}
