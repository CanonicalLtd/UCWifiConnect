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

	"gopkg.in/yaml.v2"
)

type netplanType struct {
	Network networkType `yaml:"network"`
}

type networkType struct {
	Ethernets map[string]interface{} `yaml:"ethernets,omitempty"`
	Wifis     map[string]interface{} `yaml:"wifis,omitempty"`
	Version   int                    `yaml:"version"`
}

const netplanConfig = "/etc/netplan/00-snapd-config.yaml"

func main() {

	var netplan netplanType
	buf, err1 := ioutil.ReadFile(netplanConfig)
	if err1 != nil {
		fmt.Printf("== wifi-connect: Error opening %s: %v\n", netplanConfig, err1)
		return
	}

	if err := yaml.Unmarshal(buf, &netplan); err != nil {
		fmt.Printf("== wifi-connect: Error unmarshalling yaml: %v\n", err)
		return
	}

	if netplan.Network.Wifis != nil {
		netplan.Network.Wifis = make(map[string]interface{})
		netplan.Network.Wifis["renderer"] = "NetworkManager"
	} else {
		return
	}

	out, err := yaml.Marshal(netplan)
	if err != nil {
		fmt.Printf("== wifi-connect: Error mashalling yaml: %v\n", err)
		return
	}
	fmt.Println(string(out))
	errW := ioutil.WriteFile(netplanConfig, out, 0644)
	if errW != nil {
		fmt.Printf("== wifi-connect: Error writing %s to file: %v\n", netplanConfig, errW)
	}
}
