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

	IsolatedEnvironment struct {
		UserIsolated UserIsolated `json:"user_isolated"`
		CGroup       CGroup       `json:"c_group"`

		Envs       []string `json:"envs"`
		RootfsName string   `json:"rootfs_name"`
		RootfsSHA  string   `json:"rootfs_sha"`
	}

	Record struct {
		Pid  int    `json:"pid"`
		UUID string `json:"uuid"`
	}
)
