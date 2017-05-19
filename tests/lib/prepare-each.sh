#!/bin/sh
# . $TESTSLIB/utilities.sh

# if [ $SPREAD_REBOOT -gt 0 ] ; then

#     # after reboot, configure WLAN radio interfaces
#     modprobe mac80211_hwsim radios=2
#     wait_until_interface_is_available wlan0
#     wait_until_interface_is_available wlan1

#     # Powercycle both interface to get them back into a sane state before
#     # we install the wifi-ap snap
#     snap install --devmode wireless-tools
#     for d in wlan0 wlan1 ; do
#         phy=$(iw dev $d info | awk '/wiphy/{print $2}')
#         /snap/bin/wireless-tools.rfkill block $phy
#         /snap/bin/wireless-tools.rfkill unblock $phy
#     done
#     snap remove wireless-tools

#     # first time install snaps, connect interfaces and reboot
#     # install_snap_under_test
#     # # Give wifi-ap a bit time to settle down to avoid clashed
#     # sleep 5

#     # wifi-connect.netplan

#     # # For debugging dump all snaps and connected slots/plugs
#     # snap list
#     # snap interfaces

#     # ifconfig

# else
#     # first time install snaps, connect interfaces and reboot
#     install_snap_under_test
#     # Give wifi-ap a bit time to settle down to avoid clashed
#     sleep 5

#     # For debugging dump all snaps and connected slots/plugs
#     snap list
#     snap interfaces

#     ifconfig

#     REBOOT
# fi

