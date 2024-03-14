package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/sevlyar/go-daemon"
	"visualon.com/go-server-monitor/config"
	"visualon.com/go-server-monitor/db"
	"visualon.com/go-server-monitor/server"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// End all sessions
func cleanSession() error {
	var err error

	update, err := db.DB.Prepare("UPDATE sessions SET status=?, ended_at=?")
	checkError(err)

	_, err = update.Exec("end", time.Now())
	checkError(err)

	fmt.Println("Clean sessions before exiting...")
	return nil
}

var (
	signal = flag.String("s", "", `Send signal to the daemon:
  quit — graceful shutdown
  stop — fast shutdown`)
)

func main() {
	flag.Parse()
	daemon.AddCommand(daemon.StringFlag(signal, "quit"), syscall.SIGQUIT, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)
	cntxt := &daemon.Context{
		PidFileName: "/var/run/go-server-monitor.pid",
		PidFilePerm: 0644,
		LogFileName: "/var/log/go-server-monitor.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[go-server-monitor]"},
	}

	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			log.Fatalf("Unable send signal to the daemon: %s", err.Error())
		}
		daemon.SendCommands(d)
		return
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Println("- - - - - - - - - - - - - - -")
	log.Println("daemon started")

	// Clean sessions before starting
	checkError(cleanSession())

	// Start server
	checkError(server.StartHTTPServer(config.CONFIG.Server.Port))

	err = daemon.ServeSignals()
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	log.Println("daemon terminated")
}

var (
	stop = make(chan struct{})
)

func termHandler(sig os.Signal) error {
	log.Println("terminating...")
	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		os.Exit(-1)
	}
	return daemon.ErrStop
}
