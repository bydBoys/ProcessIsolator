package impl

import (
	"ProcessIsolator/constants"
	"ProcessIsolator/impl/app/cgroups/subsystems"
	"ProcessIsolator/impl/config"
	"ProcessIsolator/impl/server/internal"
	"ProcessIsolator/util"
	"fmt"
)

type ProcessServerImpl struct {
	config.IProcessServer
	errorChan chan<- error
	msgChan   chan<- string
}

func (impl *ProcessServerImpl) Init(errorChan chan<- error, msgChan chan<- string) {
	impl.errorChan = errorChan
	impl.msgChan = msgChan
}

func (impl *ProcessServerImpl) StartProc(request *config.StartProcRequest, response *config.StartProcResponse) error {
	isUploaded := internal.CheckFile(request.IsolatedEnvironment.RootfsName)
	if !isUploaded {
		impl.errorChan <- fmt.Errorf("somebody want to StartProc, but file %s not found", request.IsolatedEnvironment.RootfsName)
		genStartProcResponse(response, "-1", fmt.Sprintf("file %s not found", request.IsolatedEnvironment.RootfsName))
		return nil
	}
	var uuid = util.RandStringBytes(10)
	process, writePipe, err := internal.MakeIsolatedProcess(&request.IsolatedEnvironment, uuid)
	if err != nil {
		impl.errorChan <- fmt.Errorf("somebody want to StartProc, but occour error %s", err)
		genStartProcResponse(response, "-1", fmt.Sprintf("internal error %s", err))
		return nil
	}

	// todo: mount
	var runPath = fmt.Sprintf(constants.ProcRunPath, uuid)
	if err = util.MountBind(runPath); err != nil {
		impl.errorChan <- fmt.Errorf("somebody want to StartProc, but occour error %s", err)
		genStartProcResponse(response, "-1", fmt.Sprintf("internal error %s", err))
		return nil
	}

	if err = process.Start(); err != nil {
		impl.errorChan <- fmt.Errorf("somebody want to StartProc, but occour error %s", err)
		genStartProcResponse(response, "-1", fmt.Sprintf("internal error %s", err))
		return nil
	}
	impl.msgChan <- fmt.Sprintf("isolatedProcess(%s) has started", uuid)
	var cgroup subsystems.ResourceConfig
	if request.IsolatedEnvironment.CGroup.Enable {
		cgroup.MemoryLimit = request.IsolatedEnvironment.CGroup.MemoryLimit
		cgroup.CpuShare = request.IsolatedEnvironment.CGroup.CpuShare
		cgroup.CpuSet = request.IsolatedEnvironment.CGroup.CpuSet
	}
	internal.PutRecord(process, uuid, &cgroup)

	util.SendCommand(request.Command, writePipe)
	impl.msgChan <- fmt.Sprintf("isolatedProcess(%s) has sended command", uuid)
	genStartProcResponse(response, uuid, "")
	return nil
}

func (impl *ProcessServerImpl) GetProcLog(request *config.GetProcLogRequest, response *config.GetProcLogResponse) error {
	lines, err := util.ReadFileLines(fmt.Sprintf(constants.ProcLogPath, request.UUID))
	if err != nil {
		impl.errorChan <- fmt.Errorf("somebody want to GetProcLog, but occour error %s", err)
		genGetProcLogResponse(response, false, nil, err.Error())
		return nil
	}
	record, err := internal.GetRecord(request.UUID)
	exist := err != nil && record != nil
	impl.msgChan <- fmt.Sprintf("isolatedProcess(%s) pid: %d status: %t", record.UUID, record.Pid, exist)
	genGetProcLogResponse(response, exist, lines, "")
	return nil
}

func (impl *ProcessServerImpl) KillProc(request *config.KillProcLogRequest, response *config.KillProcLogResponse) error {
	record, err := internal.GetRecord(request.UUID)
	exist := err != nil && record != nil
	if exist {
		if err = util.KillProcess(record.Pid); err != nil {
			impl.errorChan <- fmt.Errorf("somebody want to KillProc, but occour error %s", err)
			genKillProcResponse(response, false, err.Error())
			return nil
		}
		genKillProcResponse(response, true, "")
		return nil
	}
	impl.msgChan <- fmt.Sprintf("isolatedProcess(%s) pid: %d status: %t", record.UUID, record.Pid, false)
	genKillProcResponse(response, true, "")
	return nil
}

func (impl *ProcessServerImpl) GetVersion(request *config.GetVersionRequest, response *config.GetVersionResponse) error {
	impl.msgChan <- fmt.Sprintf("%s getVersion", request.Requester)
	*response = config.GetVersionResponse{
		Version: constants.Version,
	}
	return nil
}

// ----------------------------------------------------------------------------------------------------------

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
