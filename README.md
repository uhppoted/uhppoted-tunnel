![build](https://github.com/uhppoted/uhppoted-tunnel/workflows/build/badge.svg)

# uhppoted-tunnel

Tunnels UDP packets between a pair of machines to enable UHPPOTE controller remote access

## Status

_In development_

Supported operating systems:
- Linux
- MacOS
- Windows
- ARM7 _(e.g. RaspberryPi)_

## Raison d'être

## Releases

| *Version* | *Description*                                                                             |
| --------- | ----------------------------------------------------------------------------------------- |
|           |                                                                                           |
|           |                                                                                           |

## Installation

Executables for all the supported operating systems are packaged in the [releases](https://github.com/uhppoted/uhppoted-tunnel/releases):

The release tarballs contain the executables for all the operating systems - OS specific tarballs with all the _uhppoted_ components can be found in [uhpppoted](https://github.com/uhppoted/uhppoted/releases) releases.

Installation is straightforward - download the archive and extract it to a directory of your choice. To install `uhppoted-httpd` as a system service:
```
   cd <uhppoted directory>
   sudo uhppoted-tunnel daemonize
```

`uhppoted-tunnel help` will list the available commands and associated options (documented below).

### Building from source

Required tools:
- [Go 1.17+](https://go.dev)
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
| [uhppoted-lib](https://github.com/uhppoted/uhppoted-lib)                | common API for external applications |
| golang.org/x/sys                                                        | (for Windows service integration)    |

## uhppoted-tunnel

Usage: ```uhppoted-tunnel <command> <options>```

Supported commands:

- `help`
- `version`
- `run`
- `console`
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

