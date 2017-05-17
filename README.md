# Ubuntu Core Wifi-Connect

Wifi-connect snap allows you to connect the device to an external wifi AP. First, it puts up an AP that you can join. Then, you can open a web page the device provides that displays external APs (by SSID), where you can select one, enter the passphrase, and connect. Disconnecting later allows you to join the device AP and use the web page again to join an external AP.

* The wifi-ap snap provides the device AP.
* The network-manager snap provides management and control of the wlan0 interface used for the AP and to connect to external APs. 

## Release: Alpha 2 (0.9)

* Raspberry pi3 with no additional wifi hardware is the only verified platform

## Issue tracking

[Issues](https://github.com/CanonicalLtd/UCWifiConnect/issues)

## Use refreshed pi3 image

After installing the latest Ubuntu Core pi3 image, run:

```bash
snap refresh
```

## Install snaps

```bash
snap install wifi-ap
snap install network-manager
snap install --edge wifi-connect
```

## Create content sharing directory for wifi-ap:control interface

```bash
sudo mkdir /var/snap/wifi-connect/common/sockets
```

Note: Currently content share interface requires a reboot after connection, as described below.

(TODO: Later we'll use an interface hook script to automatically create that directory)

## Connect interfaces

```bash
snap connect wifi-connect:control wifi-ap:control
snap connect wifi-connect:network core:network
snap connect wifi-connect:network-bind core:network-bind
snap connect wifi-connect:network-manager network-manager:service
snap connect wifi-connect:network-control core:network-control
snap connect wifi-connect:network-setup-control core:network-setup-control
```

(TODO: Configure interface auto connection.)

Note: wifi-ap and network-manager interfaces auto-connect.

# Set NetWorkManager to control all networking

Note: This is a temporary manual step before network-manager snap provides a config option for this.

Note: Depending on your environment, after this you may need to use a new IP address to connect to the device.

1. Backup the existing /etc/netplan/00-snapd-config.yaml file 

```bash
sudo mv /etc/netplan/00-snapd-config.yaml ~/
```

1. Create a new netplan config file named /etc/netplan/00-default-nm-renderer.yaml:

``bash
sudo vi /etc/netplan/00-default-nm-renderer.yaml

```

Add the following two lines:

```bash
network:
    renderer: NetworkManager
```

## Reboot

Rebooting addresses a content sharing interface issue. 

Rebooting also consolidates all networking into NetworkManager.

## Optionally configure wifi-ap SSID/passphrase

If you skip these steps, the wifi-AP put up by the device has an SSID of "Ubuntu" and is unsecure (with no passphrase). 

1. Stop wifi-connect

```bash
sudo  wifi-connect stop
```

1. Set the wlan0 interface to be unmanaged by NetworkManager

```bash
nmcli set wlan0 managed n
```

1. Set the wifi-ap AP SSID

```bash
sudo  wifi-connect ssid digit
```

1. Set the AP passphrase:

```bash
sudo  wifi-connect passphrase ubuntuubuntuubuntu
```

1. Start wifi-connect

```bash
sudo  wifi-connect start
```

## Display the AP config

```bash
sudo  wifi-connect show-ap
```

Note the DHCP range:

    dhcp.range-start: 10.0.60.2
    dhcp.range-stop: 10.0.60.199

## Join the device AP

When the device AP is up and available to you, join it.

## Open the the Management Portal web page

This portal displays external wifi APs and let's you join them.

After you connect to the device AP, you can open its http portal at the .1 IP address just before the start of the DHCP range (see previous steps) using port 8080: 

    10.0.60.1:8080

### Avahi and hostname

You can also connect to the device's web page using the device host name: 

    http://HOSTNAME.local:8080 

Where HOSTNAME is the hostname of the device when it booted. (Changing hostname with the hostname command at run time is not sufficient.) 

Note: The system trying to open the web page must support Avahi. Android systems may not, for example.

## Be patient, it takes minutes

Wifi-connect pauses at daemon startup and at various times to allow state changes to settle. For example:

* On boot and on daemon start, it takes a minute or so to determine the proper state

## Disconnect from wifi

* Use `nmcli c` to display connections.
* Use `nmcli c delete CONNECTION_NAME` to disconnect and delete. This puts the device into management mode, bringing up the AP and portal.

## Logs

Log messages are currently available in journalctl and most start with "== wifi-connect", so view the system state and other messages with:

    sudo journalctl -f | grep ==

### Sample (filtered) log  

This log snippet shows the wifi-connect daemon starting, entering management mode, obtaining external SSIDs, at which point the management ap and portal are put up:

    May 05 19:06:18 localhost.localdomain snap[5990]: == wifi-connect: daemon STARTING
    May 05 19:07:06 localhost.localdomain snap[5990]: == wifi-connect: entering MANAGEMENT mode
    May 05 19:07:08 localhost.localdomain snap[5990]: == wifi-connect: SSID(s) obtained
    May 05 19:07:09 localhost.localdomain snap[5990]: == wifi-connect: start wifi-ap

The daemon waits silently until the user uses the portal to attempt to join an external AP, and, on success, the device enters operational mode:

    May 05 19:08:41 localhost.localdomain snap[5990]: 2017/05/05 19:08:41 == wifi-connect: Connecting to my_ap...
    May 05 19:08:48 localhost.localdomain snap[5990]: == wifi-connect: entering OPERATIONAL mode

Now we are connected. Verify with `nmcli c`, and then delete the connection:

    May 05 19:08:58 localhost.localdomain snap[5990]: == wifi-connect: entering MANAGEMENT mode
    May 05 19:09:00 localhost.localdomain snap[5990]: == wifi-connect: SSID(s) obtained
    May 05 19:09:02 localhost.localdomain snap[5990]: == wifi-connect: start wifi-ap

Ready to join another AP.

## Known Limitations Alpha 1

* Raspberry Pi3 with no additional hardware is the only verified platform currently 
* To set up pi3 to use wifi in console-conf, you have to reboot after first boot and run 'sudo console-conf'.
* After connecting to external wifi-ap, ifconfig shows for wlan0 the IP of the hosted AP (10.0.60.1), not the IP assigned by the external AP. But, the IP assigned by the external AP is the one that works.
* Wifi-connect takes over management of device wifi (via wlan0 interface). Any external operations that modify these may result in an incorrect state and may interrupt connectivity. For example, manually changing the network manager managed state of wlan0, or manually bringing up or down wifi-ap may break connectivity. 
* Opening the AP portal web page using device hostname (http://[hostname].local:8080) can result in a connection error from some platforms including some Android mobile phones and, in general, when connecting from any device on which Avahi is not enabled. You can open the web page using the device IP address on its AP and wlan0 interface, as described above.

## Development Environment

### Install snapd and snapcraft

```bash
sudo apt install snapd snapcraft
```

### Verify snapcraft is installed ok by printing out the version

```bash
snapcraft -v
2.27.1
```
Should output current version. More information on [snapcraft.io](https://snapcraft.io)

### Install Go

Use normal methods appropriate for your environment. 
See also: [Install Go](https://golang.org/doc/install).

### Install web development environment

Install the [Node Version Manager](https://github.com/creationix/nvm) (NVM) that will allow a specific
version of Node.js to be installed. Follow the installation instructions.

### Install the latest stable Node.js and npm

The latest stable (LTS) version of Node can be found on the [Node website](nodejs.org).

```bash
# Overview of available commands
nvm help

# Install the latest stable version
nvm install v4.4.3

# Select the version to use
nvm ls
nvm use v4.4.3
```

* Install the nodejs dependencies

```bash
npm install
```
* Update css

In case you need to update css, as gulp.js is used in this project, you would need to install it in case you haven't done that previously

```bash
npm install -g gulp
```
and execute sass task

```bash
gulp sass
```

# Pausing the daemon loop

The daemon loop can be paused with:

```bash
sudo  wifi-connect stop
```

After this, the daemon loops and does nothing. In this state you may want to run "hidden" commands (see the sourcefor these), for example to execute functions for development and verification.

Note: It is possible to execute commands that put the system into a non-working state. For example, bringing the AP UP/DOWN while wlan0 interface is managed by netork manager may result in an unworkable situation, possibly requiring reboot, or merely daemon restart.


Restart the daemon normal loop cleanly with:


```bash
sudo  wifi-connect start
```

# Unit Tests

You can run all tests by executing 
	go test -v ./...
or
	./run-check --unit


In order to run specific package test, you can:
	go test -v ./<package>

To run a specific test:
	go test -v -run <testname> 

For example:
	go test -v ./wifiap
	go test -v ./wifiap -run TestShow

More info in https://golang.org/pkg/testing

