# UCWifiConnect

Wifi-connect snap allows you to connect the device to an external wifi AP. First, it puts up an AP that you can connect to. Once connected, you can access a web portal that displays external APs (by SSID), where you can select one, enter the passphrase, and connect. 

## Release: Alpha 1

See Known Limitations below.

* Currently alpha 1 status (wifi-connect 0.6)
* Raspberry pi3 with no additional wifi hardware only tested platform
* First boot console-conf must configure and connect ethernet and NOT wifi 

### Install snaps

```bash
snap install wifi-ap
snap install network-manager
snap install --edge wifi-connect
```

### Create content sharing dir for wifi-ap:control interface

```bash
sudo mkdir /var/snap/wifi-connect/common/sockets
```

(TODO: Solution will use interface hook script when it is available to automatically create that dir)

### Connect interfaces

```bash
snap connect wifi-connect:control wifi-ap:control
snap connect wifi-connect:network core:network
snap connect wifi-connect:network-bind core:network-bind
snap connect wifi-connect:network-manager network-manager:service
snap connect wifi-connect:network-control core:network-control
```

(TODO: Configure interface auto connection.)

Note: wifi-ap and network-manager interfaces auto-connect.

Note: The content sharing interface has a known issue. Until that is resolved, You need to restart the system at this point.

### SSH to the device (ethernet) to configure AP 

(Later there may be a portal for this.)

### Stop the daemon

    sudo systemctl stop snap.wifi-connect.daemon.service

### Bring the AP down:

    sudo wifi-connect.wifi-ap -ap-off

### Set the wifi-ap AP SSID

    sudo wifi-connect.wifi-ap -ssid digit

### Set the AP passphrase:

    sudo wifi-connect.wifi-ap -passphrase ubuntuubuntuubuntu

### Start the deamon

    sudo systemctl start snap.wifi-connect.daemon.service

### Display the ap config

    sudo wifi-connect.wifi-ap -show

Note the dhcp range:

    dhcp.range-start: 10.0.60.2
    dhcp.range-stop: 10.0.60.199

After you connect to the device AP, you can open its http portal at the .1 IP address just before the start of the DCHP range: 

    10.0.60.1:8080

You can also connect to the device's AP using the machine name this way: 

    http://[hostname].local:8080 

Where [hostname] is the hostname of the device. It is a known issue that from some devices not having enabled avahi service it is not possible accessing this way (see [Limitations](#limitations) section)

## Be patient, it takes minutes

To ensure a smooth transition between states wifi-connect pauses to provide time for state changes to settle. For example:

* On boot and on daemon start, it takes a couple minutes to determine the proper state (which you can see in the log)
* When transitioning between modes (for example when connectixoign to an external AP from the web page, it takes a couple minutes  

## Logs

Log messages are currently available in journalctl and most start with "==", so view the system state and other messages with:

    sudo journalctl -f | grep ==

## Normal operations (after configuration steps)

The daemon monitors whether there's a connection to an external wifi AP using network-manager. (On start of the daemon, it waits 45 seconds to let any previous connection come up.) 

### No external AP connected

* The device is in "Management Mode" 
* Get external SSIDs until found
* The wifi-ap is put UP
* You join it
* Open the Management portal web page at IP:8080. From here you can see external APs (SSIDs), pick one and initiate the join. This takes down the current AP, the device tries to join the external AP 

### External AP is connected

* Device is in "Operational Mode"
* Daemon loops until there is no external AP connectiion known by network-manager, which causes device to be in Management Mode

Note: You can drop from external network-manager AP connections (and return the device to Management Mode) with:

    wifi-connect.netman -disconnect-wifi

(This command may  be dropped later in favor of nmcli and/or a web page.)

## Known Limitations Alpha 1

* Raspberry Pi3 with no additional hardware is the only verified platform currently 
* The device must have been configured during first boot to set up ethernet and not wifi
* Wifi-connect takes over management of the device's wlan0 interface and the wifi-ap AP. Any external operations that modify these may result in an incorrect state and may interrupt connectivity. For example, mannually changing the network manager managed state of wlan0, or manually bringing up or down wifi-ap may break connectivity. 
* Opening the AP portal web page using device hostname (http://[hostname].local:8080) can result in a connection error from some platforms including some Android mobile phones and, in general, wheni connecting from any device on which avahi is not enabled. You can open the web page using the device IP address on its AP and wlan0 interface, as described above.

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

