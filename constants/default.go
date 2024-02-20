package constants

const (
	Name    = "ProcessIsolator"
	Version = "24.2.20"
	Desc    = "A tool used for running specified process in isolated environment. Enjoy it, just for fun."

	Port        = "0.0.0.0:9963"
	BasePath    = "/run/" + Name
	ProcLogPath = BasePath + "/logs/%s.log"
	ProcRunPath = BasePath + "/runs/%s"
	FilePath    = BasePath + "/files/%s"
)
