// +build !windows

package main

import (
	"log"
	"net"
	"net/http"
	"fmt"
	"time"

	"git.zxq.co/ripple/schiavolib"
	"git.zxq.co/ripple/rippleapi/common"
	"github.com/gin-gonic/gin"
	"github.com/rcrowley/goagain"
)

func startuato(engine *gin.Engine) {
	conf, _ := common.Load()
	// Inherit a net.Listener from our parent process or listen anew.
	l, err := goagain.Listener()
	if nil != err {

		// Listen on a TCP or a UNIX domain socket (TCP here).
		if conf.Unix {
			l, err = net.Listen("unix", conf.ListenTo)
		} else {
			l, err = net.Listen("tcp", conf.ListenTo)
		}
		if nil != err {
			schiavo.Bunker.Send(err.Error())
			log.Fatalln(err)
		}

		schiavo.Bunker.Send(fmt.Sprint("LISTENINGU STARTUATO ON ", l.Addr()))

		// Accept connections in a new goroutine.
		go http.Serve(l, engine)

	} else {

		// Resume accepting connections in a new goroutine.
		schiavo.Bunker.Send(fmt.Sprint("LISTENINGU RESUMINGU ON ", l.Addr()))
		go http.Serve(l, engine)

		// Kill the parent, now that the child has started successfully.
		if err := goagain.Kill(); nil != err {
			schiavo.Bunker.Send(err.Error())
			log.Fatalln(err)
		}

	}

	// Block the main goroutine awaiting signals.
	if _, err := goagain.Wait(l); nil != err {
		schiavo.Bunker.Send(err.Error())
		log.Fatalln(err)
	}

	// Do whatever's necessary to ensure a graceful exit like waiting for
	// goroutines to terminate or a channel to become closed.
	//
	// In this case, we'll simply stop listening and wait one second.
	if err := l.Close(); nil != err {
		schiavo.Bunker.Send(err.Error())
		log.Fatalln(err)
	}
	if err := db.Close(); err != nil {
		schiavo.Bunker.Send(err.Error())
		log.Fatalln(err)
	}
	time.Sleep(time.Second * 1)
}
