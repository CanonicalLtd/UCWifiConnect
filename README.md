# Ubuntu Core Wifi-Connect

Wifi-connect snap allows you to connect the device to an external wifi AP. First, it puts up an AP that you can join. Then, you can open a web page the device provides that displays external APs (by SSID), where you can select one, enter the passphrase, and connect. Disconnecting later allows you to join the device AP and use the web page again to join an external AP.

* The wifi-ap snap provides the device AP.
* The network-manager snap provides management and control of the wlan0 interface used for the AP and to connect to external APs. 

## Release: Alpha 1

See Known Limitations below.

* Currently alpha 1 status (wifi-connect 0.6)
* Raspberry pi3 with no additional wifi hardware is the only verified platform
* At first boot, console-conf must configure and connect ethernet and NOT wifi

## Use refreshed pi3 image

After installing the latest Ubuntu Core pi3 image, run

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

(TODO: Solution will later use an interface hook script to automatically create that directory)

## Connect interfaces

```bash
snap connect wifi-connect:control wifi-ap:control
snap connect wifi-connect:network core:network
snap connect wifi-connect:network-bind core:network-bind
snap connect wifi-connect:network-manager network-manager:service
snap connect wifi-connect:network-control core:network-control
```

(TODO: Configure interface auto connection.)

Note: wifi-ap and network-manager interfaces auto-connect.

## Reboot

The content sharing interface has a known issue. Until that is resolved, you need to restart the system at this point.

## Optionally configure wifi-ap SSID/passphrase

If you skip these steps, the wifi-AP put up by the device has an SSID of "Ubuntu" and is unsecure (with no passphrase). 

1. Stop the daemon

```bash
sudo systemctl stop snap.wifi-connect.daemon.service
```

1. Bring the AP down:

```bash
sudo wifi-connect.wifi-ap -ap-off
```

1. Set the wifi-ap AP SSID

```bash
sudo wifi-connect.wifi-ap -ssid digit
```

1. Set the AP passphrase:

```bash
sudo wifi-connect.wifi-ap -passphrase ubuntuubuntuubuntu
```

1. Start the daemon

```bash
sudo systemctl start snap.wifi-connect.daemon.service
```

## Display the AP config

```bash
sudo wifi-connect.wifi-ap -show
```

Note the dhcp range:

    dhcp.range-start: 10.0.60.2
    dhcp.range-stop: 10.0.60.199

## Join the device AP

When the device AP is up and available to you, join it.

## Open the the Management Portal web page

This portal displays external wifi APs and let's you join them.

After you connect to the device AP, you can open its http portal at the .1 IP address just before the start of the DHCP range using port 8080: 

    10.0.60.1:8080

### Avahi and hostname

You can also connect to the device's web page using the device host name: 

    http://[hostname].local:8080 

Where [hostname] is the hostname of the device when it booted. (Changing hostname with the hostname command at run time is not sufficient.) 

Note: The system trying to open the web page must support Avahi. Android systems may not, for example.

## Be patient, it takes minutes

Wifi-connect pauses at daemon startup and at various times to allow state changes to settle. For example:

* On boot and on daemon start, it takes a couple minutes to determine the proper state
* When transitioning between modes (for example when connecting to an external AP from the web page, it takes a couple minutes  

## Logs

Log messages are currently available in journalctl and most start with "==", so view the system state and other messages with:

    sudo journalctl -f | grep ==

### Sample (filtered) log  

This log snippet shows the wifi-connect daemon starting (Initiation), entering Management mode, getting external SSIDs, starting the device wifi AP, and starting the Management portal, at which point it loops until the user uses the portal to attempt to join an external AP:

    Apr 28 16:34:24 localhost.localdomain snap[1766]: ======== Initiation Mode (daemon starting)
    Apr 28 16:35:50 localhost.localdomain snap[1766]: ====== Management Mode
    Apr 28 16:35:50 localhost.localdomain snap[1766]: ==== Setting wlan0 unmanaged
    Apr 28 16:35:55 localhost.localdomain snap[1766]: ==== Wifi-ap enabled state: false
    Apr 28 16:35:55 localhost.localdomain snap[1766]: ==== Setting wlan0 managed
    Apr 28 16:36:05 localhost.localdomain snap[1766]: ==== SSID(s) found and written to  /var/snap/wifi-connect/common/ssids
    Apr 28 16:36:05 localhost.localdomain snap[1766]: ==== Setting wlan0 unmanaged
    Apr 28 16:36:10 localhost.localdomain snap[1766]: ==== Have SSIDs: start wifi-ap
    Apr 28 16:36:10 localhost.localdomain snap[1766]: ==== Start Management portal if not running
    Apr 28 16:36:10 localhost.localdomain snap[1766]: ==== Writing wait file: /var/snap/wifi-connect/common/startingApConnect

This log snippet shows what happens when you select an external AP ("myap"), enter the correct password, Connecting to myap, and enter Operational Mode, where it loops until there is not external wifi connection, at which point it enters Management mode.

    Apr 28 16:37:40 localhost.localdomain snap[1766]: 2017/04/28 16:37:40 == Connecting to myap...
    Apr 28 16:39:08 localhost.localdomain snap[1766]: ======== Operational Mode
    Apr 28 16:39:08 localhost.localdomain snap[1766]: ==== Stop Management Mode http server if running

## Known Limitations Alpha 1

* Raspberry Pi3 with no additional hardware is the only verified platform currently 
* The device must have been configured during first boot to set up ethernet and not wifi
* Wifi-connect takes over management of the device's wlan0 interface and the wifi-ap AP. Any external operations that modify these may result in an incorrect state and may interrupt connectivity. For example, manually changing the network manager managed state of wlan0, or manually bringing up or down wifi-ap may break connectivity. 
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

