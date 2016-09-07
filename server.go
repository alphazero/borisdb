//    Copyright © 2016 Joubin Houshyar. All rights reserved.
//
//    This file is part of Frankinstore.
//
//    Frankinstore is free software: you can redistribute it and/or modify
//    it under the terms of the GNU Affero General Public License as
//    published by the Free Software Foundation, either version 3 of
//    the License, or (at your option) any later version.
//
//    Frankinstore is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU Affero General Public License for more details.
//
//    You should have received a copy of the GNU Affero General Public
//    License along with Frankinstore.  If not, see <http://www.gnu.org/licenses/>.

// command line tools and main server
package main

import (
	"flag"
	"fmt"
	"frankinstore/store"
	"frankinstore/web"
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
	port:   web.DefaultPort,
	dbname: store.DefaultDb,
}

/// main server process ///////////////////////////////////////////////////////

func main() {
	flag.Parse()

	// verify options
	if e := initOptions(); e != nil {
		log.Fatalf(e.Error())
	}
	log.Printf("info - frankinstore startup ... ")

	// open store
	db, e := store.OpenDb(option.path)
	if e != nil {
		log.Printf("err - failed to open database - %s", e)
		os.Exit(1)
	}
	defer db.Close()
	log.Printf("info - frankinstore using db: %q", option.path)

	// shutdown hooks
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	shutdown, shutdownFn := getShutdownHooks()

	// start webserver
	go web.RunService(option.port, db, shutdownFn)
	log.Printf("info - frankinstore listening on port %d", option.port)

	// clean shutdown
	select {
	case <-sigchan:
	case <-shutdown:
	}

	log.Printf("info - frankinstore stopped. ciao!\n")
}

/// server shutdown ///////////////////////////////////////////////////////////

func getShutdownHooks() (chan struct{}, func() error) {
	var shutdown = make(chan struct{}, 1)
	shutdownFn := func() error {
		shutdown <- struct{}{}
		close(shutdown)
		return nil
	}
	return shutdown, shutdownFn
}

/// server initialization /////////////////////////////////////////////////////

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
