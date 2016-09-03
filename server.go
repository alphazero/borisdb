package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
)

// server configuration and options
var option = struct {
	port   int    // service port
	path   string // fs store path
	dbname string // database name
}{
	port:   5722,    // default port
	dbname: "boris", // default database name
}

func init() {
	flag.IntVar(&option.port, "port", option.port, "web service port")
	flag.StringVar(&option.path, "path", option.path, "db file path")
	flag.StringVar(&option.dbname, "db", option.dbname, "db name")
}

// initialize and verify server options
func initOptions() error {
	// TODO: verify port is valid for userspace range

	// use current working directory if path is not specified.
	if option.path == "" {
		option.path = os.TempDir()
	}
	// verify db path
	finfo, e := os.Stat(option.path)
	if e != nil {
		return fmt.Errorf("err - path option - %s", e)
	}
	if !finfo.IsDir() {
		return fmt.Errorf("err - specified path is not a directory - %s", option.path)
	}

	// verify dbname
	if option.dbname == "" {
		return fmt.Errorf("err - dbname can not be blank")
	}

	option.path = filepath.Join(option.path, option.dbname)

	return nil
}

func main() {
	flag.Parse()

	// verify options
	if e := initOptions(); e != nil {
		log.Fatalf(e.Error())
	}
	log.Printf("info - frankinstore startup - db: %q", option.path)

	// TODO open store

	// TODO start webserver

	// clean shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan

	// TODO close database

	log.Printf("info - frankinstore stopped. ciao!\n")
}
