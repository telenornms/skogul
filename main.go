/*
 * gollector, main method/init
 *
 * Copyright (c) 2019 Telenor Norge AS
 * Author(s):
 *  - Kristian Lyngstøl <kly@kly.no>
 *
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2.1 of the License, or (at your option) any later version.
 * 
 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 * 
 * You should have received a copy of the GNU Lesser General Public
 * License along with this library; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA
 * 02110-1301  USA
 */

package main

import (
	"net/http"
	"time"
        "io"
	"fmt"
	"log"
	"bytes"
	"encoding/json"
)

type myHandler struct {

}

func (handler myHandler) Send(c *GollectorContainer) error {
	var buffer bytes.Buffer
	for _, m := range c.Metrics {
		fmt.Fprintf(&buffer,"test")
		for key,value := range c.Template.Metadata {
			fmt.Fprintf(&buffer,",%s=%#v",key,value)
		}
		for key,value := range m.Metadata {
			fmt.Fprintf(&buffer,",%s=%#v",key,value)
		}
		fmt.Fprintf(&buffer," ")
		comma := ""
		for key,value := range m.Data {
			fmt.Fprintf(&buffer,"%s%s=%#v",comma,key,value)
			comma = ","
		}
                lt := c.Template.Time
                if m.Time != (time.Time{}) {
                    lt = m.Time
                }
		fmt.Fprintf(&buffer," %d\n",lt.UnixNano())
	}
        log.Print("Starting backend request")
	req, err := http.NewRequest("POST", "http://127.0.0.1:8086/write?db=test", &buffer)
	req.Header.Set("Content-Type", "text/plain")
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Print(resp)
	}
        log.Print("Done")
	return nil
}

func (handler myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > 0 {
                log.Printf("Processing request from %v",r.RemoteAddr)
		b := make([]byte, r.ContentLength)
                n,err := io.ReadFull(r.Body,b)
                if err != nil {
                    log.Panicf("Read error from client, read %d bytes: %s", n,err)
                }
		var m GollectorContainer
		err = json.Unmarshal(b,&m)
		if err == nil {
			err = m.Validate()
		}
		if err == nil {
			handler.Send(&m)
			fmt.Fprintf(w, "OK")
		} else {
			fmt.Fprintf(w, "Unable to parse JSON: %s", err)
		}
                log.Printf("Done with %v",r.RemoteAddr)
	}
}

func main() {
	http.Handle("/", myHandler{})
	log.Fatal(http.ListenAndServe(":8080",nil))
}
