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

package netplan

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	netplanConfig    = "/etc/netplan/00-snapd-config.yaml"
	netplanOutConfig = netplanConfig
)

type Client struct {
	NetplanFile    string
	NetplanOutFile string
}

type NetplanType struct {
	Network NetworkType `yaml:"network"`
}

type NetworkType struct {
	Ethernets map[string]interface{} `yaml:"ethernets,omitempty"`
	Wifis     map[string]interface{} `yaml:"wifis,omitempty"`
	Version   int                    `yaml:"version"`
}

// the runtime client object
func RuntimeClient() *Client {
	client := &Client{}
	client.NetplanFile = netplanConfig
	client.NetplanOutFile = netplanOutConfig
	return client
}

// set the config file
func (c *Client) SetConfigFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	c.NetplanFile = path
	c.NetplanOutFile = path
	return nil
}

// set the config file to write to (useful in testing)
func (c *Client) SetConfigOutFile(path string) {
	c.NetplanOutFile = path
	return
}

// Open netplan config fail, parse yaml, and return data
func (c *Client) GetNetplan() (*NetplanType, error) {

	netplan := &NetplanType{}

	buf, err1 := ioutil.ReadFile(c.NetplanFile)
	if err1 != nil {
		fmt.Printf("== wifi-connect/netplan: Error opening %s: %v\n", netplanConfig, err1)
		return netplan, err1
	}

	if err := yaml.Unmarshal(buf, &netplan); err != nil {
		fmt.Printf("== wifi-connect/netplan: Error unmarshalling yaml: %v\n", err)
		return netplan, err
	}
	return netplan, nil
}

// Modify netplan config data to set wifis to be managed by network manager
func (c *Client) SetWifisNetworkManager(config *NetplanType) {

	if config.Network.Wifis != nil {
		config.Network.Wifis = make(map[string]interface{})
		config.Network.Wifis["renderer"] = "NetworkManager"
	} else {
		fmt.Println("== wifi-connect/netplan. The netplan config does not manage 'wifis', no change to it is needed")
	}
	return
}
