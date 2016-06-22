// +build windows

package main

import (
	"net"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"git.zxq.co/ripple/rippleapi/common"
)

func startuato(engine *gin.Engine) {
	conf, _ := common.Load()
	var (
		l net.Listener
		err error
	)
	// Listen on a TCP or a UNIX domain socket (TCP here).
	if conf.Unix {
		l, err = net.Listen("unix", conf.ListenTo)
	} else {
		l, err = net.Listen("tcp", conf.ListenTo)
	}
	if nil != err {
		log.Fatalln(err)
	}

	http.Serve(l, engine)
}
