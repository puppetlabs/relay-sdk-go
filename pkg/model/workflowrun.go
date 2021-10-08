package model

type WorkflowRun struct {
	Name      string `json:"name"`
	RunNumber int32  `json:"run_number"`
	URL       string `json:"url"`
}
