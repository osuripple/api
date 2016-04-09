package v1

import (
	"fmt"
	"os"
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

var upSince time.Time

// MetaUpSinceGET retrieves the moment the API application was started.
// Mainly used to get if the API was restarted.
func MetaUpSinceGET(md common.MethodData) common.Response {
	return common.Response{
		Code: 200,
		Data: upSince.UnixNano(),
	}
}
