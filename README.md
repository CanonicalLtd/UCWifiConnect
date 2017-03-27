# UCWifiConnect
Snap to be able to switch a wireless card of a UC device into AP mode and use it to configure wireless

## Development Environment

### Install snapd and snapcraft
```bash
sudo apt install snapd snapcraft

# Verify snapcraft is installed ok by printing out the version
snapcraft -v
2.27.1
```
should output current version. More information on [snapcraft.io](https://snapcraft.io)

### Install Go
Follow the instructions to [install Go](https://golang.org/doc/install).

- Install NVM
Install the [Node Version Manager](https://github.com/creationix/nvm) that will allow a specific
version of Node.js to be installed. Follow the installation instructions.

- Install the latest stable Node.js and npm
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

- Install the nodejs dependencies
```bash
cd identity-vault
npm install
```
