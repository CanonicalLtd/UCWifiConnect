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

package daemonPkg

import (
	"os"
	"testing"

	"github.com/CanonicalLtd/UCWifiConnect/utils"
)

func TestManualFlagPath(t *testing.T) {
	mfpInit := "mfp"
	client := GetClient()
	client.SetManualFlagPath(mfpInit)
	mfp := client.GetManualFlagPath()
	if mfp != mfpInit {
		t.Errorf("ManualFlag path should be %s but is %s\n", mfpInit, mfp)
	}
}

func TestWaitFlagPath(t *testing.T) {
	wfpInit := "wfp"
	client := GetClient()
	client.SetWaitFlagPath(wfpInit)
	wfp := client.GetWaitFlagPath()
	if wfp != wfpInit {
		t.Errorf("WaitFlag path should be %s but is %s", wfp, wfpInit)
	}
}

func TestState(t *testing.T) {
	client := GetClient()
	client.SetState(MANAGING)
	client.SetPreviousState(STARTING)
	ps := client.GetPreviousState()
	if ps != STARTING {
		t.Errorf("Previous state should be %d but is %d", STARTING, ps)
	}
	s := client.GetState()
	if s != MANAGING {
		t.Errorf("State should be %d but is %d", MANAGING, s)
	}
}

func TestCheckWaitApConnect(t *testing.T) {
	client := GetClient()
	wfp := "thispathnevershouldexist"
	client.SetWaitFlagPath(wfp)
	if client.CheckWaitApConnect() {
		t.Errorf("CheckWaitApConnect returns true but should return false")
	}
	wfp = "../static/tests/waitFile"
	client.SetWaitFlagPath(wfp)
	if !client.CheckWaitApConnect() {
		t.Errorf("CheckWaitApConnect returns false but should return true")
	}
}

func TestManualMode(t *testing.T) {
	mfp := "thisfileshouldneverexist"
	client := GetClient()
	client.SetManualFlagPath(mfp)
	client.SetState(MANUAL)
	if client.ManualMode() {
		t.Errorf("ManualMode returns true but should return false")
	}
	if client.GetState() != STARTING {
		t.Errorf("ManualMode should set state to STARTING when not in manual mode but does not")
	}
	mfp = "../static/tests/manualMode"
	client.SetManualFlagPath(mfp)
	client.SetState(STARTING)
	if !client.ManualMode() {
		t.Errorf("ManualMode returns false but should return true")
	}
	if client.GetState() != MANUAL {
		t.Errorf("ManualMode should set state to MANUAL when in manual mode but does not")
	}
}

func TestSetDefaults(t *testing.T) {
	client := GetClient()
	hfp := "/tmp/hash"
	utils.SetHashFile(hfp)
	client.SetDefaults()
	_, err := os.Stat(utils.HashFile)
	if os.IsNotExist(err) {
		t.Errorf("SetDefaults should have created %s but did not", hfp)
	}
	res, _ := utils.MatchingHash("wifi-connect")
	if !res {
		t.Errorf("SetDefaults password match did not match")
	}
}
