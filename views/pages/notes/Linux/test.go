package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

const cgroupMemoryHierarchMount = "/sys/fs/cgroup/memory"

func main() {
	cmd := exec.Command("sh", "-c", `stress --vm-bytes 200m --vm-keep -m 1`)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Printf("%v\n", cmd.Process.Pid)
		os.Mkdir(path.Join(cgroupMemoryHierarchMount, "testmemorylimit"), 0755)
		ioutil.WriteFile(path.Join(cgroupMemoryHierarchMount, "testmemorylimit", "tasks"), []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
		ioutil.WriteFile(path.Join(cgroupMemoryHierarchMount, "testmemorylimit", "memory.limit_in_bytes"), []byte("50m"), 0644)
	}
	cmd.Process.Wait()
}
