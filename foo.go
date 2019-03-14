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

type GollectorContainer struct {
	Source	string			`json:"src"`
	Metrics	[]GollectorMetric	`json:"metrics"`
}

type GollectorMetric struct {
	Time		time.Time	`json:"timestamp"`
	Metadata	map[string]interface{}	`json:"metadata"`
	Data		map[string]interface{}	`json:"data"`
}

type gerror struct {
	Reason string
}

func (e gerror) Error() string {
	return e.Reason
}

type myHandler struct {

}

func (m GollectorMetric) validate() error {
	if m.Data == nil {
		return gerror{"Missing data for metric"}
	}
	/*
	 * This is a bit of a cheat, due to a lack of a decent nil-test
	 */
	if m.Time == (time.Time{}) {
		return gerror{"no time?"}
	}
	return nil
}
func (c GollectorContainer) validate() error {
	if c.Source != "test" {
		return gerror{"Only support \"test\" data source for now"}
	}
	if c.Metrics == nil {
		return gerror{"Missing metrics[] data"}
	}
	if len(c.Metrics) <= 0 {
		return gerror{"Empty metrics[] data"}
	}
	for i := 0; i < len(c.Metrics); i++ {
		err := c.Metrics[i].validate()
		if err != nil {
			return err
		}
	}
	return nil;
}

func (handler myHandler) Send(c *GollectorContainer) error {
	var buffer bytes.Buffer
	for _, m := range c.Metrics {
		fmt.Fprintf(&buffer,"%s",c.Source)
		for key,value := range m.Metadata {
			fmt.Fprintf(&buffer,",%s=%#v",key,value)
		}
		fmt.Fprintf(&buffer," ")
		comma := ""
		for key,value := range m.Data {
			fmt.Fprintf(&buffer,"%s%s=%#v",comma,key,value)
			comma = ","
		}
		fmt.Fprintf(&buffer," %d\n",m.Time.UnixNano())
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
			err = m.validate()
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
