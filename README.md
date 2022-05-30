![build](https://github.com/uhppoted/uhppoted-tunnel/workflows/build/badge.svg)

# uhppoted-tunnel

Tunnels UDP packets between a pair of machines to enable UHPPOTE controller remote access.

Technically it's not really a tunnel, except in the sense that as a packet you enter a dark
forbidding hole, mysterious and possibly unspeakable things occur and you emerge some time
later blinking in the light in an entirely different place. So probably more a relay or a 
proxy .. but we're going to go with _tunnel_ anyway.

## Raison d'être

For those **so** annoying times when it would ne nice to run the UHPPOTE _AccessControl_ application
but the controller is in one place and the host machine is in another (or perhaps on a VPS in Norway) 
which means UDP broadcast doesn't just work. And poking holes in the firewall and/or fixing the NAT
or setting up a VPN is just more trouble than it's worth.

Also useful for remotely using
- [uhppote-cli](https://github.com/uhppoted/uhppote-cli)
- [uhppoted-app-sheets](https://github.com/uhppoted/uhppoted-app-sheets) 
- [uhppoted-app-wildapricot](https://github.com/uhppoted/uhppoted-app-wildapricot)

and is a simpler alternative to:
- [uhppoted-rest](https://github.com/uhppoted/uhppoted-rest)
- [uhppoted-mqtt](https://github.com/uhppoted/uhppoted-mqtt)

## Status

_In development_

Supported operating systems:
- Linux
- MacOS
- Windows
- ARM7 _(e.g. RaspberryPi)_

## Releases

| *Version* | *Description*                                                                             |
| --------- | ----------------------------------------------------------------------------------------- |
|           |                                                                                           |
|           |                                                                                           |

## Installation

Executables for all the supported operating systems are packaged in the [releases](https://github.com/uhppoted/uhppoted-tunnel/releases):

The release tarballs contain the executables for all the operating systems - OS specific tarballs with all the _uhppoted_ components can be found in [uhpppoted](https://github.com/uhppoted/uhppoted/releases) releases.

Installation is straightforward - download the archive and extract it to a directory of your choice. To install `uhppoted-tunnel` as a system service:
```
   cd <uhppoted directory>
   sudo uhppoted-tunnel daemonize
```

`uhppoted-tunnel help` will list the available commands and associated options (documented below).

### Building from source

Required tools:
- [Go 1.18+](https://go.dev)
- make (optional but recommended)

To build using the included Makefile:

```
git clone https://github.com/uhppoted/uhppoted-tunnel.git
cd uhppoted-tunnel
make build
```

Without using `make`:
```
git clone https://github.com/uhppoted/uhppoted-tunnel.git
cd uhppoted-tunnel
go build -trimpath -o bin ./...
```

The above commands build the `uhppoted-tunnel` executable to the `bin` directory.


#### Dependencies

| *Dependency*                                                            | *Description*                        |
| ----------------------------------------------------------------------- | -------------------------------------|
| [uhppote-core](https://github.com/uhppoted/uhppote-core)                | Device level API implementation      |
| [uhppoted-lib](https://github.com/uhppoted/uhppoted-lib)                | Common library functions             |
| golang.org/x/sys                                                        | (for Windows service integration)    |

## uhppoted-tunnel

Usage: ```uhppoted-tunnel <command> <options>```

Supported commands:

- `help`
- `version`
- `run`
- `daemonize`
- `undaemonize`

Defaults to `run` if the command it not provided i.e. ```uhppoted-tunnel <options>``` is equivalent to 
```uhppoted-tunnel run <options>```.

### `run`

Runs the `uhppoted-tunnel` service. Default command, intended for use as a system service that runs in the 
background. 

Command line:

` uhppoted-tunnel [--debug] [--console]`

```
  --console     Runs the HTTP server endpoint as a console application, logging events to the console.
  --debug       Displays verbose debugging information, in particular the communications with the 
                UHPPOTE controllers
```

### `daemonize`

Registers `uhppoted-tunnel` as a system service that will be started on system boot. The command creates the necessary
system specific service configuration files and service manager entries. On Linux it defaults to using the 
`uhppoted:uhppoted` user:group - this can be changed with the `--user` option

Command line:

`uhppoted-tunnel daemonize [--user <user>]`

### `undaemonize`

Unregisters `uhppoted-tunnel` as a system service, but does not delete any created log or configuration files. 

Command line:

`uhppoted-tunnel undaemonize `

## Supporting files

