package app

import (
	"ProcessIsolator/util"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func RunIsolatedProcessInit() error {
	cmdArray := util.ReadCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("isolated process get command error")
	}

	setUpMount()

	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		return err
	}
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		return err
	}
	return nil
}

func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println("Get current location error ", err)
		return
	}

	if err := pivotRoot(pwd); err != nil {
		log.Println("pivot ", err)
		return
	}

	// 挂载 proc
	if err := syscall.Mount("proc", "/proc", "proc", syscall.MS_NOEXEC|syscall.MS_NOSUID|syscall.MS_NODEV, ""); err != nil {
		log.Println("proc ", err)
		return
	}
	// 挂载 devtmpfs
	//25 28 0:6 / /dev rw,nosuid,relatime shared:2 - devtmpfs udev rw,size=491380k,nr_inodes=122845,mode=755
	err = syscall.Mount("devtmpfs", "/dev", "devtmpfs", syscall.MS_NOSUID|syscall.MS_RELATIME, "rw,size=491380k,mode=755")
	if err != nil {
		log.Println("dev ", err)
	}
	// 挂载 devpts
	//26 25 0:23 / /dev/pts rw,nosuid,noexec,relatime shared:3 - devpts devpts rw,gid=5,mode=620,ptmxmode=000
	err = syscall.Mount("devpts", "/dev/pts", "devpts", syscall.MS_NOEXEC|syscall.MS_NOSUID|syscall.MS_RELATIME, "rw,mode=620")
	if err != nil {
		log.Println("devpts ", err)
	}
	// 挂载 devshm
	//30 25 0:25 / /dev/shm rw,nosuid,nodev shared:4 - tmpfs tmpfs rw
	err = syscall.Mount("tmpfs", "/dev/shm", "tmpfs", syscall.MS_NOSUID|syscall.MS_NODEV, "rw")
	if err != nil {
		log.Println("devshm ", err)
	}
}

/*
p.s. chroot只改变当前进程的root目录，而pivotRoot会改变当前ns的root路径
1. 将当前路径（容器的root）通过bind制作为挂载点
2. 在当前路径下创建临时路径.pivot_root
3. 通过PivotRoot系统调用，将当前路径设为新root路径，并将曾经的路径挂载到.pivot_root上
4. 调用chdir，将"/"设置为新工作路径
5. 解挂临时路径.pivot_root
*/
func pivotRoot(currentPath string) error {
	if err := syscall.Mount(currentPath, currentPath, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount rootfs to itself error: %v", err)
	}
	pivotDir := filepath.Join(currentPath, ".pivot_root")
	if err := os.MkdirAll(pivotDir, 0777); err != nil {
		return fmt.Errorf("mkdir %v", err)
	}
	if err := syscall.PivotRoot(currentPath, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}
	if err := os.Remove(pivotDir); err != nil {
		return fmt.Errorf("remove pivot_root dir %v", err)
	}
	return nil
}
