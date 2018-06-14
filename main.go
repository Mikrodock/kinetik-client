// Example of a daemon with echo service
package main

import (
	"fmt"
	"kinetik-client/agent"
	"kinetik-server/logger"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/subosito/gotenv"

	"github.com/takama/daemon"
)

const (

	// name of the service
	name        = "kinetik-client"
	description = "Agent and executor for Mikrodock"

	// port which daemon should be listen
	port = ":9977"
)

//    dependencies that are NOT required by the service, but might be used
var dependencies = []string{"dummy.service"}

var stdlog, errlog *log.Logger

// Service has embedded daemon
type Service struct {
	daemon.Daemon
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {

	usage := "Usage: kinetik-client install | remove | start | stop | status"

	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}

	agent.StartSampling()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	for {
		select {
		case killSignal := <-interrupt:
			if killSignal == os.Interrupt {
				return "Kinetik-Client was interruped by system signal", nil
			}
			return "Kinetik-Client was killed", nil
		}
	}

	return "Kinetik-Client exiting", nil

}

func init() {

	gotenv.Load("/root/.env")

	logFile, _ := os.OpenFile("/var/log/kinetik.out", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	errFile, _ := os.OpenFile("/var/log/kinetik.err", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	stdlog = log.New(logFile, "", log.Ldate|log.Ltime)
	errlog = log.New(errFile, "", log.Ldate|log.Ltime)
	logger.StdLog = stdlog
	logger.ErrLog = errlog
}

func main() {
	srv, err := daemon.New(name, description, dependencies...)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
