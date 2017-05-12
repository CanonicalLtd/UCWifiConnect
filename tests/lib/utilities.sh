#!/bin/sh

wait_for_systemd_service() {
	while ! systemctl status $1 ; do
		sleep 1
	done
	sleep 1
}

wait_for_systemd_service_exit() {
	while systemctl status $1 ; do
		sleep 1
	done
	sleep 1
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
		# Install prebuilt snap
		snap install --devmode ${PROJECT_PATH}/${SNAP_NAME}_*_${SNAP_ARCH}.snap
		# Setup all necessary aliases
		snapd_version=$(snap version | awk '/^snapd / {print $2; exit}')
		for alias in $SNAP_AUTO_ALIASES ; do
			target=$SNAP_NAME.$alias
			if dpkg --compare-versions $snapd_version lt 2.25 ; then
				target=$SNAP_NAME
			fi
			snap alias $target $alias
		done
	fi
}

install_snap() {
	# Don't reinstall if we have it installed already
	if ! snap list | grep $1 ; then
		snap install --$2 $1
	fi
}

install_additional_snaps() {
	install_snap wifi-ap stable
	install_snap network-manager stable
}

connect_interfaces() {
	snap connect wifi-ap:network-manager network-manager:service

	snap connect wifi-connect:control wifi-ap:control
	snap connect wifi-connect:network-manager network-manager:service
	snap connect wifi-connect:network
	snap connect wifi-connect:network-bind
	snap connect wifi-connect:network-control
	snap connect wifi-connect:network-setup-control
}

