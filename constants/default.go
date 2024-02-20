package constants

const (
	Name    = "ProcessIsolator"
	Version = "24.2.20"
	Desc    = "A tool used for running specified process in isolated environment. Enjoy it, just for fun."

	Port        = "0.0.0.0:9963"
	ProcLogPath = "/run/" + Name + "/logs/%s.log"
	ProcRunPath = "/run/" + Name + "/runs/%s"
	FilePath    = "/run/" + Name + "/files/%s"
)
