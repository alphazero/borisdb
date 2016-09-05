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
	"strings"
)

var option = struct {
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

var data string

func init() {
	flag.StringVar(&data, "d", data, "data to send")
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
	run(client)
}

func run(client *web.Client) {

	resp, e := client.Put([]byte(data))
	if e != nil {
		fmt.Printf("err - %s - ", e)
	}
	resp = strings.Trim(resp, " \n")
	fmt.Printf("[%s]\n", resp)
}
