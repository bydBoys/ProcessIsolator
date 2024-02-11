package config

type (
	Record struct {
		Pid  int    `json:"pid"`
		UUID string `json:"uuid"`
	}

	StartProcRequest struct {
		Commands     []string `json:"commands"`
		UserIsolated `json:"user_isolated"`
		CGroup       `json:"c_group"`
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

	IProcServer interface {
		StartProc(request *StartProcRequest, response *StartProcResponse) error
		GetProcLog(request *GetProcLogRequest, response *GetProcLogResponse) error
	}
)
