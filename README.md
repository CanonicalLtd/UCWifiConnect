# UCWifiConnect
wifi-connect snap allows you to connect the device to an external wifi AP. First, it puts up an AP that you can connect to. Once connected, you can access a web portal that displays external APs (by SSID), where you can select one, enter the passphrase, and connect. 

## Status

* Currently alpha status (wifi-connect 0.4)
* Works on pi3 with no additional wifi hardware

## Set up

### Install snaps

```bash
snap install wifi-ap
snap install network-manager
snap install --edge wifi-connect
```

### Create content sharing dir for wifi-ap:control interfaces
```bash
sudo mkdir /var/snap/wifi-connect/common/sockets
```

(TODO: Solution will use interface hook script to create that dir the first time)

### Connect interfaces

```bash
snap connect wifi-connect:control wifi-ap:control
snap connect wifi-connect:network core:network
snap connect wifi-connect:network-bind core:network-bind
snap connect wifi-connect:network-manager network-manager:service
snap connect wifi-connect:network-control core:network-control
```

(TODO configure auto connection)

Note: wifi-ap and network-manager interfaces should auto-connect.

Note: The content sharing interface has a known bug. You may need to restart the system.

### SSH to the device (ethernet) to configure AP 

(Later there will be a portal for this.)

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

After you connect to the device's AP, you can open its http portal at 

    10.0.60.1:8080

you can also connect to the device's AP using the machine name this way: 

    http://[hostname].local:8080 

where [hostname] is the hostname of the device. It is a known issue that from some devices not having enabled avahi service it is not possible accessing this way (see [Limitations](#limitations) section)

## Normal operations (after configuration steps)

The daemon monitors whether there's a connection to an external wifi AP. (On start of the daemon, it waits 45 seconds to let any previous connection come up.) 

### No external AP connected

* The device is in Management Mode
* Get external SSIDs until found
* The wifi-ap is put UP
* You join it
* Open the Management portal web page at IP:8080. From here you can see external APs (SSIDs), pick one and initiate the join. This takes down the current AP, the device tries to join the external AP 

### External AP is connected

* Device is in Operational Mode
* Operational Port is put UP (this has no content now, but will allow setting the device to Management Mode later)
* Connect to Operational port via AP IP address (TODO: avahi)
* Daemon loops until there is no external AP connectiion, which causes device to be in Managerment Mode

Note: Until we have an Operational Portal, you can drop from external AP connections with:

    wifi-connect.netman -disconnect-wifi

## Various commands

Running commands may interfere with normal operations controlled by the daemon. Before running any commands stop the daemon (explained above), then restart it after.

### Network Manager dbus commands 

```bash
wifi-connect.netman -help
Usage of netman:
  -check-connected
        Check if connected at all
  -check-connected-wifi
        Check if connected to external wifi
  -disconnect-wifi
        Disconnect from any and all external wifi
  -get-ssids
        Only display SSIDs (don't connect)
  -manage-iface string
        Set the specified interface to be managed by network-manager.
  -unmanage-iface string
        Set the specified interface to NOT be managed by network-manager.
  -wifis-managed
        Show list of wifi interfaces that are managed by network-manager
```

### Wifi-ap commands

```bash
sudo wifi-connect.wifi-ap -help
Usage of wifi-ap:
  -ap-off
        Turn off the AP
  -ap-on
        Turn on the AP
  -enabled
        Check if the AP is UP
  -passphrase string
        Set the AP's passphrase
  -show
        Show the wifi-ap configuration
  -ssid string
        Set the AP's SSID
  -verbose
        Display verbose output
```

Note: On version 0.4 if you Turn off the AP, you also need to manually clear a state file:

    sudo rm /var/snap/wifi-connect/common/startingApConnect

### Additional

* Most log messages start with "==" for viewing with $sudo journalctl -f | grep ==
* display current wifi status on device with: nmcli d

## Development Environment

### Install snapd and snapcraft
```bash
sudo apt install snapd snapcraft

### Verify snapcraft is installed ok by printing out the version
snapcraft -v
2.27.1
```
should output current version. More information on [snapcraft.io](https://snapcraft.io)

### Install Go
Follow the instructions to [install Go](https://golang.org/doc/install).

### Install web development environment

- Install NVM
Install the [Node Version Manager](https://github.com/creationix/nvm) that will allow a specific
version of Node.js to be installed. Follow the installation instructions.

### Install the latest stable Node.js and npm
The latest stable (LTS) version of Node can be found on the [Node website](nodejs.org).
```bash
# Overview of available commands
nvm help

### Install the latest stable version
nvm install v4.4.3

### Select the version to use
nvm ls
nvm use v4.4.3
```

- Install the nodejs dependencies
```bash
npm install
```

- Update css
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

to run a specific test:
	go test -v -run <testname> 

for example:
	go test -v ./wifiap
	go test -v ./wifiap -run TestShow

More info in https://golang.org/pkg/testing


# Limitations

When accession AP portal in browser using device hostname (http://[hostname].local:8080) could result in a connection error. This
is something known when accessing from some Android mobile phones and, in general, if connecting from a not avahi enabled device
