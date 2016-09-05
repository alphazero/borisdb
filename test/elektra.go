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

package main

import (
	"flag"
	"fmt"
	"frankinstore/web"
	"os"
	"strings"
)

var option = struct {
	cmd   string
	data  string
	host  string
	port  int
	size  int
	count int
}{
	host:  "127.0.0.1",
	port:  web.DefaultPort,
	size:  4096,
	count: 1,
}

func init() {
	flag.StringVar(&option.cmd, "c", option.cmd, "command: {put, get, shutdown, info}")
	flag.StringVar(&option.data, "d", option.data, "data to send")
	flag.StringVar(&option.host, "a", option.host, "host address")
	flag.IntVar(&option.port, "p", option.port, "port")
	flag.IntVar(&option.size, "s", option.size, "size of payload")
	flag.IntVar(&option.count, "n", option.count, "number of concurrent requests")
}

func main() {
	flag.Parse()

	client, e := web.NewClient(option.host, option.port)
	if e != nil {
		fmt.Printf("err - %s\n", e)
		return
	}
	switch option.cmd {
	case "put":
		put(client)
	case "get":
		get(client)
	case "info":
	case "shutdown":
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", option.cmd)
		os.Exit(1)
	}
}

func put(client *web.Client) {

	resp, e := client.Put([]byte(option.data))
	if e != nil {
		fmt.Printf("err - %s - ", e)
	}
	resp = strings.Trim(resp, " \n")
	fmt.Printf("[%s]\n", resp)
}

func get(client *web.Client) {

	resp, e := client.Get(option.data)
	if e != nil {
		fmt.Printf("err - %s - ", e)
	}
	rs := strings.Trim(string(resp), " \n")
	fmt.Printf("[%s]\n", rs)
}
