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

package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

// PrintMapSorted prints to stdout a map sorting content by keys
func PrintMapSorted(m map[string]interface{}) {
	sortedKeys := make([]string, 0, len(m))
	for key := range m {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)
	for _, k := range sortedKeys {
		fmt.Fprintf(os.Stdout, "%s: %v\n", k, m[k])
	}
}

// WriteWaitFile writes the "wait" file used as a flag by the daemon
func WriteWaitFile() {
	wait := os.Getenv("SNAP_COMMON") + "/startingApConnect"
	err := ioutil.WriteFile(wait, []byte("wait"), 0644)
	if err != nil {
		fmt.Println("== wifi-connect: Error writing flag wait file:", err)
	}
}

// RemoveWaitFile removes the "wait" file used as a flag by the daemon
func RemoveWaitFile() {
	waitApPath := os.Getenv("SNAP_COMMON") + "/startingApConnect"
	if _, err := os.Stat(waitApPath); !os.IsNotExist(err) {
		err := os.Remove(waitApPath)
		if err != nil {
			//try again after pause
			time.Sleep(0000 * time.Millisecond)
			os.Remove(waitApPath)
		}
	}
}
