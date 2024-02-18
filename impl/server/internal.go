package server

import (
	"ProcessIsolator/constants"
	"ProcessIsolator/impl/app/cgroups"
	"ProcessIsolator/impl/app/cgroups/subsystems"
	"ProcessIsolator/impl/config"
	"ProcessIsolator/util"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
)

// 可以并发读，但写时需互斥
var rwLock sync.RWMutex
var records = make(map[string]config.Record)

func waitProc(c *exec.Cmd, id string, cgroupManager *cgroups.CgroupManager) {
	runtime.Gosched()
	_ = c.Wait()
	rwLock.Lock()
	delete(records, id)
	rwLock.Unlock()
	runtime.Gosched()
	cgroupManager.Destroy()
}

func putRecord(cmd *exec.Cmd, uuid string, cgroup *subsystems.ResourceConfig) config.Record {
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
	return record
}

func getRecord(uuid string) (*config.Record, error) {
	rwLock.RLock()
	record, exist := records[uuid]
	rwLock.RUnlock()
	if !exist {
		return nil, fmt.Errorf("UUID %s not exist", uuid)
	}
	return &record, nil
}

func newIsolatedProcess(isolatedEnv config.IsolatedEnvironment, uuid string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := util.NewPipe()
	if err != nil {
		return nil, nil
	}
	initCmd, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return nil, nil
	}

	path, err := getRootPath(uuid, isolatedEnv.RootfsName, isolatedEnv.RootfsSHA)
	if err != nil {
		return nil, nil
	}

	cmd := exec.Command(initCmd, "init")
	var cloneFlags uintptr = syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWTIME
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   cloneFlags,
		Unshareflags: syscall.CLONE_NEWNS,
	}
	if isolatedEnv.UserIsolated.Enable {
		cloneFlags |= syscall.CLONE_NEWUSER
		cmd.SysProcAttr.Credential = &syscall.Credential{
			Uid: uint32(1),
			Gid: uint32(1),
		}
	}
	logPath := fmt.Sprintf(constants.ProcLogPath, uuid)

	if cmd.Stdout, err = util.GenerateFile(logPath); err != nil {
		log.Fatal("generate logFile error", err)
		return nil, nil
	}

	cmd.ExtraFiles = []*os.File{readPipe}
	cmd.Env = isolatedEnv.Envs
	cmd.Dir = path

	return cmd, writePipe
}
