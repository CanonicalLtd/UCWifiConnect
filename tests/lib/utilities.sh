#!/bin/sh

wait_for_systemd_service() {
	while ! systemctl status $1 ; do
		sleep 1
	done
	sleep 1
}

wait_for_daemon_ready() {
	wait_for_systemd_service snap.wifi-connect.daemon.service
}

wait_for_systemd_service_exit() {
	count=20
	while systemctl status $1 && count -gt 0; do
		sleep 1
		let count--
	done
	sleep 1

	if [ count -eq 0 ]; then
		exit 1
	fi
}

does_interface_exist() {
	[ -d /sys/class/net/$1 ]
}

wait_until_interface_is_available() {
	while ! does_interface_exist $1; do
		# Wait for 200ms
		sleep 0.2
	done
}

stop_after_first_reboot() {
	if [ $SPREAD_REBOOT -gt 0 ] ; then
		exit 0
	fi
}

install_snap() {
	# Don't reinstall if we have it installed already
	if ! snap list | grep $1 ; then
		snap install --$2 $1
	fi
}

install_additional_snaps() {
	install_snap network-manager stable
	install_snap wifi-ap stable
}

connect_interfaces() {
	#snap connect wifi-connect:control wifi-ap:control
	snap connect wifi-connect:network core:network
	snap connect wifi-connect:network-bind core:network-bind
	snap connect wifi-connect:network-manager network-manager:service
	snap connect wifi-connect:network-control core:network-control
	snap connect wifi-connect:network-setup-control core:network-setup-control
}

install_snap_under_test() {
	# If we don't install the snap here we get a system
	# without any network connectivity after reboot.
	if [ -n "$SNAP_CHANNEL" ] ; then
		# Don't reinstall if we have it installed already
		if ! snap list | grep $SNAP_NAME ; then
			snap install --$SNAP_CHANNEL $SNAP_NAME
		fi
	else
		
		
		install_additional_snaps

		# Install prebuilt snap
		snap install --devmode ${PROJECT_PATH}/${SNAP_NAME}_*_${SNAP_ARCH}.snap

		# Create content sharing directory if needed
		[ -e /var/snap/wifi-connect/common/sockets ] || mkdir -p /var/snap/wifi-connect/common/sockets

		connect_interfaces

		# Setup all necessary aliases
		snapd_version=$(snap version | awk '/^snapd / {print $2; exit}')
		for alias in $SNAP_AUTO_ALIASES ; do
			target=$SNAP_NAME.$alias
			if dpkg --compare-versions $snapd_version lt 2.25 ; then
				target=$SNAP_NAME
			fi
			[ -e /snap/bin/$alias ] || snap alias $target $alias
		done
	fi
}

