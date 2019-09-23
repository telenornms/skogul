/*
 * skogul, using config file
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

/*
skogul-file parses a json-based config file and starts skogul.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/KristianLyng/skogul/config"
	"github.com/KristianLyng/skogul/sender"
)

var ffile = flag.String("file", "~/.config/skogul.json", "Path to skogul config to read.")
var frecvhelp = flag.String("receiver-help", "", "Print extra options for receiver")
var fsendhelp = flag.String("sender-help", "", "Print extra options for sender")
var ftarget = flag.String("sender", "debug://", "Where to send data. See -help for details.")
var fhelp = flag.Bool("help", false, "Print extensive help/usage")

func help() {
	flag.Usage()
}

func helpSender(s string) {
	if sender.Auto[s] == nil {
		fmt.Printf("No such sender %s\n", s)
		return
	}
	sh, _ := config.HelpSender(s)
	sh.Print()
}

func main() {
	flag.Parse()
	if *fhelp {
		help()
		os.Exit(0)
	}
	if *fsendhelp != "" {
		helpSender(*fsendhelp)
		os.Exit(0)
	}

	_, err := config.File(*ffile)
	if err != nil {
		log.Fatal(err)
	}

}
