package config

type (
	Record struct {
		Pid  int    `json:"pid"`
		UUID string `json:"uuid"`
	}

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

	IProcServer interface {
		StartProc(request *StartProcRequest, response *StartProcResponse) error
		GetProcLog(request *GetProcLogRequest, response *GetProcLogResponse) error
		KillProc(request *KillProcLogRequest, response *KillProcLogResponse) error
		GetVersion(request *int, response *string) error
	}
)
