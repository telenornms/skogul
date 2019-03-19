/*
 * skogul, http writer
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

package senders

import (
	"bytes"
	"encoding/json"
	"github.com/KristianLyng/skogul/pkg"
	"log"
	"net/http"
	"time"
)

type HTTP struct {
	URL string
}

func (ht HTTP) Send(c *skogul.Container) error {
	b, err := json.Marshal(*c)
	var buffer bytes.Buffer
	buffer.Write(b)
	req, err := http.NewRequest("POST", ht.URL, &buffer)
	req.Header.Set("Content-Type", "application/json")
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return err
	} else {
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			log.Print(resp)
		}
	}
	return nil
}
