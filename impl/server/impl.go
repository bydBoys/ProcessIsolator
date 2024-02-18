package server

import (
	"ProcessIsolator/constants"
	"ProcessIsolator/impl/app/cgroups/subsystems"
	"ProcessIsolator/impl/config"
	"ProcessIsolator/util"
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
func genKillProcResponse(response *config.KillProcLogResponse, success bool, errMsg string) {
	*response = config.KillProcLogResponse{
		Success: success,
		Error:   errMsg,
	}
}

func (receiver ProcServerImpl) StartProc(request *config.StartProcRequest, response *config.StartProcResponse) error {
	var uuid = util.RandStringBytes(10)
	child, writePipe := newIsolatedProcess(request.IsolatedEnvironment, uuid)
	if child == nil {
		genStartProcResponse(response, "-1", "fork child error")
		return nil
	}
	if err := child.Start(); err != nil {
		genStartProcResponse(response, "-1", "start child error"+err.Error())
		return nil
	}
	var cgroup subsystems.ResourceConfig
	if request.IsolatedEnvironment.CGroup.Enable {
		cgroup.MemoryLimit = request.IsolatedEnvironment.CGroup.MemoryLimit
		cgroup.CpuShare = request.IsolatedEnvironment.CGroup.CpuShare
		cgroup.CpuSet = request.IsolatedEnvironment.CGroup.CpuSet
	}
	putRecord(child, uuid, &cgroup)
	genStartProcResponse(response, uuid, "")

	util.SendCommand(request.Command, writePipe)
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

func (receiver ProcServerImpl) KillProc(request *config.KillProcLogRequest, response *config.KillProcLogResponse) error {
	record, err := getRecord(request.UUID)
	exist := err != nil && record != nil
	if exist {
		if err := util.KillProcess(record.Pid); err != nil {
			genKillProcResponse(response, false, err.Error())
			return nil
		}
		genKillProcResponse(response, true, "")
		return nil
	}
	genKillProcResponse(response, true, "")
	return nil
}

func (receiver ProcServerImpl) GetVersion(request *int, response *string) error {
	*response = constants.Version
	return nil
}
