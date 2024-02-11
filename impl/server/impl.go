package server

import (
	"ProcZygote/constants"
	"ProcZygote/impl/app/cgroups/subsystems"
	"ProcZygote/impl/config"
	"ProcZygote/util"
	"fmt"
)

type ProcServerImpl struct {
	config.IProcServer
}

func genStartProcResponse(response *config.StartProcResponse, uuid string, errMsg string) {
	*response = config.StartProcResponse{
		Error: errMsg,
		UUID:  uuid,
	}
}
func genGetProcLogResponse(response *config.GetProcLogResponse, exist bool, logs []string, errMsg string) {
	*response = config.GetProcLogResponse{
		Exist: exist,
		Error: errMsg,
		Logs:  logs,
	}
}

func (receiver ProcServerImpl) StartProc(request *config.StartProcRequest, response *config.StartProcResponse) error {
	var uuid = util.RandStringBytes(10)
	child, writePipe := newChildProcess(request.UserIsolated, constants.Hook, uuid)
	if child == nil {
		genStartProcResponse(response, "-1", "fork child error")
		return nil
	}
	if err := child.Start(); err != nil {
		genStartProcResponse(response, "-1", "start child error"+err.Error())
		return nil
	}
	var cgroup subsystems.ResourceConfig
	if request.CGroup.Enable {
		cgroup.MemoryLimit = request.CGroup.MemoryLimit
		cgroup.CpuShare = request.CGroup.CpuShare
		cgroup.CpuSet = request.CGroup.CpuSet
	}
	putRecord(child, uuid, &cgroup)
	genStartProcResponse(response, uuid, "")

	util.SendCommand(request.Commands, writePipe)
	return nil
}

func (receiver ProcServerImpl) GetProcLog(request *config.GetProcLogRequest, response *config.GetProcLogResponse) error {
	lines, err := util.ReadFileLines(fmt.Sprintf(constants.ProcLogPath, request.UUID))
	if err != nil {
		genGetProcLogResponse(response, false, nil, err.Error())
		return nil
	}
	record, err := getRecord(request.UUID)
	exist := err != nil && record != nil
	genGetProcLogResponse(response, exist, lines, "")
	return nil
}
