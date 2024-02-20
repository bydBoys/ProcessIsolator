package internal

import (
	"ProcessIsolator/constants"
	"ProcessIsolator/impl/app/cgroups"
	"ProcessIsolator/impl/app/cgroups/subsystems"
	"ProcessIsolator/impl/config"
	"ProcessIsolator/util"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
)

// 可以并发读，但写时需互斥
var rwLock sync.RWMutex
var records = make(map[string]config.Record)

func MakeIsolatedProcess(environment *config.IsolatedEnvironment, uuid string) (*exec.Cmd, *os.File, error) {
	readPipe, writePipe, err := util.NewPipe()
	if err != nil {
		return nil, nil, err
	}
	initCmd, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return nil, nil, err
	}

	cmd := exec.Command(initCmd, "init")
	var cloneFlags uintptr = syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWTIME
	if environment.UserIsolated.Enable {
		cloneFlags |= syscall.CLONE_NEWUSER
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   cloneFlags,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	if cmd.Stdout, err = util.GenerateFile(fmt.Sprintf(constants.ProcLogPath, uuid)); err != nil {
		return nil, nil, fmt.Errorf("generate logFile error %s", err)
	}
	var runPath = fmt.Sprintf(constants.ProcRunPath, uuid)
	if err = os.MkdirAll(fmt.Sprintf(constants.ProcRunPath, uuid), 0600); err != nil {
		return nil, nil, err
	}

	if err = util.UnTar(fmt.Sprintf(constants.FilePath, environment.RootfsName), runPath); err != nil {
		return nil, nil, err
	}

	cmd.Env = environment.Envs
	cmd.Dir = runPath

	cmd.ExtraFiles = []*os.File{readPipe}
	return cmd, writePipe, nil
}

func PutRecord(cmd *exec.Cmd, uuid string, cgroup *subsystems.ResourceConfig) {
	var record = config.Record{
		Pid:  cmd.Process.Pid,
		UUID: uuid,
	}
	rwLock.Lock()
	records[uuid] = record
	rwLock.Unlock()
	cgroupManager := cgroups.NewCgroupManager(uuid)
	if cgroup != nil {
		cgroupManager.Set(cgroup)
		cgroupManager.Apply(cmd.Process.Pid)
	}
	go waitProc(cmd, uuid, cgroupManager)
}

func GetRecord(uuid string) (*config.Record, error) {
	rwLock.RLock()
	record, exist := records[uuid]
	rwLock.RUnlock()
	if !exist {
		return nil, fmt.Errorf("UUID %s not exist", uuid)
	}
	return &record, nil
}

// ----------------------------------------------------------------------------------------------------------

func waitProc(c *exec.Cmd, id string, cgroupManager *cgroups.CgroupManager) {
	runtime.Gosched()
	_ = c.Wait()
	rwLock.Lock()
	delete(records, id)
	rwLock.Unlock()
	runtime.Gosched()
	cgroupManager.Destroy()
	// todo: unmount
}
