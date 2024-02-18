package constants

const (
	Name    = "ProcessIsolator"
	Version = "24.2.18"
	Desc    = "A tool used for running specified process in isolated environment. Enjoy it, just for fun."

	Port        = "0.0.0.0:9963"
	ProcLogPath = "/run/" + Name + "/logs/%s.log"
)
