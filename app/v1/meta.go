package v1

import (
	"os"
	"syscall"

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
	go proc.Signal(syscall.SIGUSR2)
	return
}
