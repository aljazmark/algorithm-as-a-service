package algorithms

//Algorithm interface definition for algorithm packages
type Algorithm interface {
	//Handles input parsing and algorithm execution
	run(Parameters) (string, string, error)
}

//Parameters structure definition
type Parameters struct {
	Algorithm  string   `json:"algorithm,omitempty"`
	Input      string   `json:"input,omitempty"`
	Parameters []string `json:"parameters,omitempty"`
}

//RunAlgorithm calls algorithm run method and handles recovery
func RunAlgorithm(a Algorithm, request Parameters) (string, string, error) {
	output, time, err := a.run(request)
	if err != nil {
		return "", "", err
	}
	return output, time, nil
}
