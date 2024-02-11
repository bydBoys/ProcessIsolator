package constants

import "ProcZygote/impl/config"

const (
	Name    = "ProcZygote"
	Version = "24.1.31"
	Usage   = "ProcZygote is a simple container runtime implementation.Enjoy it, just for fun."

	Port        = "0.0.0.0:9963"
	ProcLogPath = "/run/" + Name + "/logs/%s.log"
)

var Hook = func(isolated config.UserIsolated, uuid string) {

}
