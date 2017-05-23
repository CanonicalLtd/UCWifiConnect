#!/bin/bash

. $TESTSLIB/utilities.sh

if [ $SPREAD_REBOOT -gt 0 ] ; then

    # after reboot, configure WLAN radio interfaces
    modprobe mac80211_hwsim radios=2
    wait_until_interface_is_available wlan0
    wait_until_interface_is_available wlan1

    # Powercycle both interface to get them back into a sane state before
    # we install the wifi-ap snap
    snap install --devmode wireless-tools
    for d in wlan0 wlan1 ; do
        phy=$(iw dev $d info | awk '/wiphy/{print $2}')
        /snap/bin/wireless-tools.rfkill block $phy
        /snap/bin/wireless-tools.rfkill unblock $phy
    done
    snap remove wireless-tools

else
	echo "Wait for firstboot change to be ready"
	while ! snap changes | grep -q "Done"; do
		snap changes || true
		snap change 1 || true
		sleep 1
	done

	echo "Ensure fundamental snaps are still present"
	. $TESTSLIB/snap-names.sh
	for name in $gadget_name $kernel_name $core_name; do
		if ! snap list | grep -q $name ; then
			echo "Not all fundamental snaps are available, all-snap image not valid"
			echo "Currently installed snaps:"
			snap list
			exit 1
		fi
	done

	echo "Kernel has a store revision"
	snap list | grep ^${kernel_name} | grep -E " [0-9]+\s+canonical"

	# Snapshot of the current snapd state for a later restore
	if [ ! -f $SPREAD_PATH/snapd-state.tar.gz ] ; then
		systemctl stop snapd.service snapd.socket
		tar czf $SPREAD_PATH/snapd-state.tar.gz /var/lib/snapd
		systemctl start snapd.socket
	fi

    # first time install snaps, connect interfaces and reboot
    install_snap_under_test

    # For debugging dump all snaps and connected slots/plugs
    snap list
    snap interfaces

    REBOOT
fi








# For debugging dump all snaps and connected slots/plugs
snap list
snap interfaces



