package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

func main() {
	appName := os.Args[1]
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	appDir := filepath.Join(filepath.Dir(exe), "..", appName)
	pid := fmt.Sprintf("%d", os.Getpid())
	goodCgroupPath := filepath.Join("/sys/fs/cgroup/cpu/good", appName)
	badCgroupPath := filepath.Join("/sys/fs/cgroup/cpu/bad", appName)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	go func() {
		<-sigs
		must(ioutil.WriteFile("/sys/fs/cgroup/cpu/cgroup.procs", []byte(pid), 0755))
		must(deleteCgroup(goodCgroupPath))
		must(deleteCgroup(badCgroupPath))
		must(os.RemoveAll(appDir))
		os.Exit(1)
	}()

	must(os.MkdirAll(appDir, 0755))
	must(os.MkdirAll(goodCgroupPath, 0755))

	must(ioutil.WriteFile(filepath.Join(appDir, "pidfile"), []byte(pid), 0755))
	must(ioutil.WriteFile(filepath.Join(goodCgroupPath, "cgroup.procs"), []byte(pid), 0755))
	must(ioutil.WriteFile(filepath.Join(goodCgroupPath, "cpu.shares"), []byte("1000"), 0755))

	for i := 0; i < runtime.NumCPU(); i++ {
		go eatCPU(appDir)
	}

	select {}
}

func deleteCgroup(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	for {
		if err := os.Remove(path); err != nil {
			time.Sleep(100 * time.Millisecond)
		} else {
			return nil
		}
	}
}

func eatCPU(appDir string) {
	for {
		if _, err := os.Stat(filepath.Join(appDir, "spike")); os.IsNotExist(err) {
			time.Sleep(time.Second)
		} else {
			for i := 0; i < 100000000; i++ {
				// eat cpu
			}
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
