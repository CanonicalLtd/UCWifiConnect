#!/bin/bash

# Simulate a WiFi radio network interface
#modprobe mac80211_hwsim radios=2

# We don't have to build a snap when we should use one from a
# channel
if [ -n "$SNAP_CHANNEL" ] ; then
	exit 0
fi

# If there is a wifi-connect snap prebuilt for us, lets take
# that one to speed things up.
if [ -e ${PROJECT_PATH}/${SNAP_NAME}_*_${SNAP_ARCH}.snap ] ; then
	exit 0
fi

# Search for updates
#snap refresh

# Setup classic snap and build the wifi-connect snap in there
snap install --devmode --beta classic

cat <<-EOF > /home/test/build-snap.sh
#!/bin/sh
set -ex
apt update
apt install -y --force-yes snapcraft
cd ${PROJECT_PATH}
snapcraft clean
snapcraft
EOF
chmod +x /home/test/build-snap.sh
sudo classic /home/test/build-snap.sh
snap remove classic

# Make sure we have a snap build
test -e ${PROJECT_PATH}/${SNAP_NAME}_*_${SNAP_ARCH}.snap