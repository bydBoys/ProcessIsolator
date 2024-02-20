package config

type (
	StartProcRequest struct {
		IsolatedEnvironment IsolatedEnvironment
		Command             []string `json:"command"`
	}
	StartProcResponse struct {
		Error string `json:"error"`
		UUID  string `json:"uuid"`
	}

	GetProcLogRequest struct {
		UUID string `json:"uuid"`
	}
	GetProcLogResponse struct {
		Exist bool     `json:"exist"`
		Error string   `json:"error"`
		Logs  []string `json:"logs"`
	}

	KillProcLogRequest struct {
		UUID string `json:"uuid"`
	}
	KillProcLogResponse struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
	}
)
