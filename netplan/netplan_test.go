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
	"testing"
	//"gopkg.in/yaml.v2"
)

func TestSetConfigFile(t *testing.T) {
	path := "filedoesnotexist"
	c := RuntimeClient()
	err := c.SetConfigFile(path)
	if err == nil {
		t.Errorf("A yaml file %s that does not exist failed to provoke error:", path)
	}

	path = "../static/tests/00-snapd-config.yaml"
	err = c.SetConfigFile(path)
	if err != nil {
		t.Errorf("A yaml file %s that exists provoked an error %s", path, err)
	}
}

func TestGetNetplan(t *testing.T) {

	badPath := "notexist"
	out := "/tmp/out.yml"
	c := RuntimeClient()
	c.SetConfigFile(badPath)
	c.SetConfigOutFile(out)
	_, err := c.GetNetplan()
	if err == nil {
		t.Errorf("Failed to report error when opening non-existing netplan file %s", badPath, err)
	}
	path := "../static/tests/00-snapd-config.yaml"
	c.SetConfigFile(path)

	netPlan, err1 := c.GetNetplan()
	if err1 != nil {
		t.Errorf("Failed to report error when opening non-existing netplan file %s", badPath, err1)
	}
	if _, ok := netPlan.Network.Ethernets["eth0"]; !ok {
		t.Errorf("Netplan file does not contain 'eth0' but it should")
	}
	if _, ok := netPlan.Network.Wifis["wlan0"]; !ok {
		t.Errorf("Netplan file does not contain 'wlan0'et0' but it should")
	}

	fmt.Println(netPlan)
}

func TestSetWifisMetworkManager(t *testing.T) {
	path := "../static/tests/00-snapd-config.yaml"
	out := "/tmp/out.yml"
	c := RuntimeClient()
	c.SetConfigFile(path)
	c.SetConfigOutFile(out)
	np, _ := c.GetNetplan()
	fmt.Println(np)
	c.SetWifisNetworkManager(np)
	fmt.Println(np.Network.Wifis["renderer"])
	if np.Network.Wifis == nil {
		t.Errorf("iModified Netplan config 'wifis' does not contain 'renderer' key but should")
	}

	if np.Network.Wifis["renderer"] != "NetworkManager" {
		t.Errorf("Modified netplan config 'wifis'>'renderer' does not equal NetworkManager' but should")

	}
}
