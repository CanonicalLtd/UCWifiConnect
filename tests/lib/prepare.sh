#!/bin/bash

. $TESTSLIB/utilities.sh

stop_after_first_reboot

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

install_snap_under_test
install_additional_snaps

# Snapshot of the current snapd state for a later restore
if [ ! -f $SPREAD_PATH/snapd-state.tar.gz ] ; then
	systemctl stop snapd.service snapd.socket
	tar czf $SPREAD_PATH/snapd-state.tar.gz /var/lib/snapd
	systemctl start snapd.socket
fi

# Create content sharing directory
[ -e /var/snap/wifi-connect/common/sockets ] || mkdir -p /var/snap/wifi-connect/common/sockets

connect_interfaces

# For debugging dump all snaps and connected slots/plugs
snap list
snap interfaces

# netplan file needs to be modified to set wlan0 as managed by network manager
wifi-connect.netplan

REBOOT