package config

type (
	UserIsolated struct {
		Enable bool `json:"enable"`
	}
	CGroup struct {
		Enable      bool   `json:"enable"`
		CpuShare    string `json:"cpu_share"`
		CpuSet      string `json:"cpu_set"`
		MemoryLimit string `json:"memory_limit"`
	}
)
