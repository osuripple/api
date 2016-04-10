package v1

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/osuripple/api/common"
)

// MetaRestartGET restarts the API with Zero Downtime™.
func MetaRestartGET(md common.MethodData) (r common.Response) {
	proc, err := os.FindProcess(syscall.Getpid())
	if err != nil {
		r.Code = 500
		r.Message = "couldn't find process. what the fuck?"
		return
	}
	r.Code = 200
	r.Message = "brb"
	go func() {
		time.Sleep(time.Second)
		proc.Signal(syscall.SIGUSR2)
	}()
	return
}

// MetaKillGET kills the API process. NOTE TO EVERYONE: NEVER. EVER. USE IN PROD.
// Mainly created because I couldn't bother to fire up a terminal, do htop and kill the API each time.
func MetaKillGET(md common.MethodData) (r common.Response) {
	proc, err := os.FindProcess(syscall.Getpid())
	if err != nil {
		r.Code = 500
		r.Message = "couldn't find process. what the fuck?"
		return
	}
	const form = "02/01/2006"
	r.Code = 200
	// yes
	r.Message = fmt.Sprintf("RIP ripple API %s - %s", upSince.Format(form), time.Now().Format(form))
	go func() {
		time.Sleep(time.Second)
		proc.Kill()
	}()
	return
}

var upSince = time.Now()

// MetaUpSinceGET retrieves the moment the API application was started.
// Mainly used to get if the API was restarted.
func MetaUpSinceGET(md common.MethodData) common.Response {
	return common.Response{
		Code: 200,
		Data: upSince.UnixNano(),
	}
}

// MetaUpdateGET updates the API to the latest version, and restarts it.
func MetaUpdateGET(md common.MethodData) common.Response {
	if f, err := os.Stat(".git"); err == os.ErrNotExist || !f.IsDir() {
		return common.Response{
			Code:    500,
			Message: "repo is not using git",
		}
	}
	go func() {
		// go get
		//        -u: update all dependencies (including API source)
		//        -d: stop after downloading deps
		if !execCommand("go", "get", "-u", "-d") {
			return
		}
		if !execCommand("go", "build") {
			return
		}

		proc, err := os.FindProcess(syscall.Getpid())
		if err != nil {
			log.Println(err)
			return
		}
		proc.Signal(syscall.SIGUSR2)
	}()
	return common.Response{
		Code:    200,
		Message: "Started updating! " + surpriseMe(),
	}
}

func execCommand(command string, args ...string) bool {
	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return false
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)
		return false
	}
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return false
	}
	data, err := ioutil.ReadAll(stderr)
	if err != nil {
		log.Println(err)
		return false
	}
	// Bob. We got a problem.
	if len(data) != 0 {
		log.Println(string(data))
		return false
	}
	io.Copy(os.Stdout, stdout)
	cmd.Wait()
	stdout.Close()
	return true
}
