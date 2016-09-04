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
)

var option = struct {
	host  string
	port  int
	count int
}{
	host:  "127.0.0.1",
	port:  8080,
	count: 1,
}

var data string

func init() {
	flag.StringVar(&data, "k", data, "object key")
	flag.StringVar(&option.host, "a", option.host, "host address")
	flag.IntVar(&option.port, "p", option.port, "port")
	flag.IntVar(&option.count, "n", option.count, "number of concurrent requests")
}

func main() {
	fmt.Printf("Salaam!\n")
	flag.Parse()

	for i := 0; i < option.count; i++ {
		run()
	}
}

func run() {

	client, e := web.NewClient(option.host, option.port)
	if e != nil {
		fmt.Printf("err - %s\n", e)
		return
	}

	value, e := client.Get(data)
	if e != nil {
		fmt.Printf("err - %s\n", e)
		return
	}

	fmt.Printf("ok - %s\n", string(value))
}
