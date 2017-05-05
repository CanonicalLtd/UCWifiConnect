/*
 * Copyright (C) 2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/CanonicalLtd/UCWifiConnect/netplan"

	"gopkg.in/yaml.v2"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	Config    string `short:"n" long:"netplan-file" description:"Netplan config file (to override default)"`
	ConfigOut string `short:"o" long:"netplan--out-file" description:"Netplan config file to write to (to override default, which is the Config file)"`
}

func main() {
	c := netplan.RuntimeClient()

	_, err := flags.ParseArgs(&opts, os.Args[1:])
	if err != nil {
		fmt.Println("== wifi-connect/netplan. Error parsing arguments:", err)
		return
	}

	if len(opts.Config) > 1 {
		err1 := c.SetConfigFile(opts.Config)
		if err1 != nil {
			fmt.Println("== wifi-connect/netplan. Specified config file does not exist:", opts.Config)
			return
		}
	}

	if len(opts.ConfigOut) > 1 {
		c.SetConfigOutFile(opts.ConfigOut)
	}

	netplanConfig, err := c.GetNetplan()

	if err != nil {
		fmt.Println("== wifi-connect/netplan. Error getting netplan config:", err)
	}

	fmt.Println(netplanConfig)

	c.SetWifisNetworkManager(netplanConfig)

	out, err := yaml.Marshal(netplanConfig)

	if err != nil {
		fmt.Printf("== wifi-connect: Error mashalling yaml: %v\n", err)
		return
	}

	errW := ioutil.WriteFile(c.NetplanOutFile, out, 0644)
	if errW != nil {
		fmt.Printf("== wifi-connect: Error writing %s to file: %v\n", netplanConfig, errW)
		return
	}

	fmt.Println("== wifi-connect:netplan. Netplan config file", c.NetplanOutFile, "now specifies that wifis are managed by NetworkManager. Reboot or run 'sudo netplan generate && sudo netplan apply' to make the config live. **Caution**: this will break any existing network connection over wifi.")
}
