![build](https://github.com/uhppoted/uhppoted-tunnel/workflows/build/badge.svg)

# uhppoted-tunnel

Tunnels UDP packets between a pair of machines to enable UHPPOTE controller remote access.

Technically it's not really a tunnel, except in the sense that as a packet you enter a dark
forbidding hole, mysterious and possibly unspeakable things occur and you emerge some time
later blinking in the light in an entirely different place. So probably more a relay or a 
proxy .. but we're going with _tunnel_ anyway.

The implementation includes the following connectors:
- UDP listen
- UDP broadcast
- TCP server
- TCP client
- TLS server
- TLS client
- HTTP POST
- HTTPS POST

## Raison d'être

For those **so** annoying times when it would be nice to run the UHPPOTE _AccessControl_ application
but the controller is in one place and the host machine is in another (or perhaps on a VPS in Norway) 
which means UDP broadcast doesn't just work. And poking holes in the firewall and tweaking the NAT
or setting up a VPN is either not going to happen or is more trouble than it's worth.

Also useful for remotely using
- [uhppote-cli](https://github.com/uhppoted/uhppote-cli)
- [uhppoted-app-sheets](https://github.com/uhppoted/uhppoted-app-sheets) 
- [uhppoted-app-wildapricot](https://github.com/uhppoted/uhppoted-app-wildapricot)

and is a simpler alternative to:
- [uhppoted-rest](https://github.com/uhppoted/uhppoted-rest)
- [uhppoted-mqtt](https://github.com/uhppoted/uhppoted-mqtt)

## Status

Supported operating systems:
- Linux
- MacOS
- Windows
- ARM7 _(e.g. RaspberryPi)_

## Releases

| *Version* | *Description*                                                                             |
| --------- | ----------------------------------------------------------------------------------------- |
| v0.8.1    | Expanded on HTTP/S example application                                                    |
| v0.8.0    | Initial release                                                                           |

## Installation

Executables for all the supported operating systems are packaged in the [releases](https://github.com/uhppoted/uhppoted-tunnel/releases):

The release tarballs contain the executables for all the operating systems - OS specific tarballs with all the _uhppoted_ components can be found in [uhpppoted](https://github.com/uhppoted/uhppoted/releases) releases.

Installation is straightforward - download the archive and extract it to a directory of your choice. To install `uhppoted-tunnel` as a system service:
```
   cd <uhppoted directory>
   sudo uhppoted-tunnel daemonize --in <connector> --out <connector> --label <label>
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
go build -trimpath -o bin/ ./...
```

The above commands build the `uhppoted-tunnel` executable to the `bin` directory.


#### Dependencies

| *Dependency*                                                            | *Description*                        |
| ----------------------------------------------------------------------- | -------------------------------------|
| [uhppote-core](https://github.com/uhppoted/uhppote-core)                | Device level API implementation      |
| [uhppoted-lib](https://github.com/uhppoted/uhppoted-lib)                | Common library functions             |
| golang.org/x/sys                                                        | (for Windows service integration)    |

## uhppoted-tunnel

Usage: ```uhppoted-tunnel <command> --in <connector> --out <connector> <options>```

Supported commands:

- `help`
- `version`
- `run`
- `daemonize`
- `undaemonize`

Defaults to `run` if the command it not provided i.e. ```uhppoted-tunnel --in <connector> --out <connector> <options>``` is equivalent to 
```uhppoted-tunnel run  --in <connector> --out <connector> <options>```.

### `run`

Runs the `uhppoted-tunnel` service. Default command, intended for use as a system service that runs in the 
background. 

Command line:

` uhppoted-tunnel [--debug] [--console] --in <connector> --out <connector> [options]`

```
  --in <connector>  Defines the connector that accepts incoming commands. Valid 'in' connectors include: 
                    - udp/listen:<bind address> (e.g. udp/listen:0.0.0.0:60000)
                    - tcp/server:<bind address> (e.g. tcp/server:0.0.0.0:12345)
                    - tcp/client:<host address> (e.g. tcp/client:192.168.1.100:12345)
                    - tls/server:<bind address> (e.g. tls/server:0.0.0.0:12345)
                    - tls/client:<host address> (e.g. tls/client:192.168.1.100:12345)
                    - http/<bind address> (e.g. http/0.0.0.0:8080)
                    - https/<bind address> (e.g. https/0.0.0.0:8443)

  --out <connector> Defines the connector that forwards received commands. Valid 'out' connectors include: 
                    - udp/broadcast:<broadcast address> (e.g. udp/broadcast:255.255.255.255:60000)
                    - tcp/server:<bind address> (e.g. tcp/server:0.0.0.0:12345)
                    - tcp/client:<host address> (e.g. tcp/client:192.168.1.100:12345)
                    - tls/server:<bind address> (e.g. tls/server:0.0.0.0:12345)
                    - tls/client:<host address> (e.g. tls/client:192.168.1.100:12345)

  --console     Runs the UDP tunnel as a console application, logging events to the console.
  --debug       Displays verbose debugging information, in particular the communications with the 
                UHPPOTE controllers

  Options:

  --max-retries <retries>  Maximum number of failed bind/connect attempts before failing with a fatal error.
                           Defaults to 32, set to -1 for infinite retry.                

  --max-retry-delay <delay>  Retries use an exponential backoff (starting at 5 seconds) up to the delay (in
                             human readable time format e.g. 60s or 5m). Defaults to 5 minutes.

  --lockfile <file>  Overrides the default lockfile name for use in e.g. bash scripts. The default lockfile
                     name is generated from the hash of the 'in' and 'out' connectors.

  --log-level <level>  Lowest level log messages to include in logging output ('debug', 'info', 'warn' or 'error'). 
                       Defaults to 'info'
```

In general, tunnels operate in pairs - one on the _host_, listening for commands from e.g. the _AccessControl_ application
or _uhppote-cli_ and the other on the _client_ local to the controller, which sends the commands to the controller(s)
and returns the replies to the _host_. It is however, possible to chain multiple tunnels to bridge across several machines.

### `daemonize`

Registers `uhppoted-tunnel` as a system service that will be started on system boot. The command creates the necessary
system specific service configuration files and service manager entries. On Linux it defaults to using the 
`uhppoted:uhppoted` user:group - this can be changed with the `--user` option

Command line:

`uhppoted-tunnel daemonize --in <connector> --out <connector> [--label <label>] [--user <user>]`

```
  --label <label>  Identifying label for the tunnel daemon/service, used to identify the tunnel in logs and when
                   uninstalling the daemon/service. Imperative if running multiple tunnel daemons on the same machine,
                   optional but recommended otherwise. Defaults to uhppoted-tunnel if not provided.
```

### `undaemonize`

Unregisters `uhppoted-tunnel` as a system service, but does not delete any created log or configuration files. 

Command line:

`uhppoted-tunnel undaemonize `


```
  --label <label>  Identifying label for the tunnel daemon/service to be uninstalled. Defaults to uhppoted-tunnel if
                   not provided.
```

## Connectors

_uhppoted-tunnel_ includes support for multiple connectors which can in general be mixed and matched, with some restrictions:

_IN_ connectors:

- UDP listen
- TCP server
- TCP client
- TLS server
- TLS client
- HTTP POST
- HTTPS POST

_OUT_ connectors:

- UDP broadcast
- TCP server
- TCP client
- TLS server
- TLS client

### UDP listen

Listens for incoming UDP packets on the _bind address_, effectively acting as a direct proxy for a remote controller.

```
--in udp/listen:<bind address>

e.g. 

--in udp/listen:0.0.0.0:60000
```

### UDP broadcast

Sends a received packet out as a UDP message on the _broadcast address_ and forwards any replies to the original requester,
effectively acting as a proxy for a remote application.

```
--out udp/broadcast:<broadcast address> [--udp-timeout <timeout>]

   The broadcast address is typically (but not necessarily) the UDP broadcast for the network adapter for the controllers'
   network segment. However it can be any valid IPv4 address:port combination to accomodate the requirements of the 
   installation.

   --udp-timeout <timeout>  Sets the maximum time to wait for replies to a broadcast message, in human readable format
                            e.g. 15s, 1250ms, etc. Defaults to 5 seconds if not provided.

e.g. 

--out udp/broadcast:255.255.255.255:60000 --udp-timeout 5s
```

### TCP server

The TCP server connector accepts connections from one or more TCP clients and can act as both an _IN_ connector and an _OUT_ connector.
Incoming requests will be forwarded to all connected clients.

```
--in tcp/server:<bind address>

e.g. 

--in tcp/server:0.0.0.0:12345
```

### TCP client

The TCP client connector connects to a TCP server and can act as both an _IN_ connector and an _OUT_ connector. Incoming requests/replies
will be forwarded to the remote server.

```
--in tcp/client:<host address>

e.g. 

--in tcp/host:192.168.1.100:12345
```

### TLS server

The TLS server connector is a TCP server connector that only accepts TLS secured client connections.

```
--in tls/server:<bind address> [--ca-cert <file>] [--cert <file>] [--key <file>] [--client-auth]

  --ca-cert      CA certificate used to verify client certificates (defaults to ca.cert)
  --cert         server TLS certificate in PEM format (defaults to server.cert)
  --key          server TLS key in PEM format (defaults to server.key)
  --client-auth  requires client mutual authentication if supplied

e.g. 

--in tls/server:0.0.0.0:12345 --ca-cert tunnel.ca --cert tunnel.cert --key tunnel.key --client-auth
```

### TLS client

The TLS client connector is a TCP client connector that only connects to TLS secured servers.

```
--in tls/client:<host address> [--ca-cert <file>] [--cert <file>] [--key <file>] [--client-auth]

  --ca-cert      CA certificate used to verify server certificates (defaults to ca.cert)
  --cert         client TLS certificate in PEM format. Optional, only required if the TLS server 
                 has mutual authentication enabled.
  --key          client TLS key in PEM format. Optional, only required if the TLS server 
                 has mutual authentication enabled.

e.g. 

--in tls/client:192.168.1.100:12345 --ca-cert tunnel.ca --cert client.cert --key client.key
```

### HTTP POST

The HTTP POST connector accepts JSON POST requests and forwards replies to the requesting client, primarily
to support quick and dirty browser based applications (a tiny example is included in the _examples_ folder).

```
--in http/<bind address> [--html <folder>]

  --html <folder> Folder containing the HTML served to the browser on the bind address.

e.g. 

--in http:/0.0.0.0:8080 --html examples/html
```

POST request:
```
  {
    ID: <request ID>,
    wait: <UDP timeout>,
    request: <UDP request byte array>
  }

e.g.

  {
    ID: 19,
    wait: "5s",
    request: [0x17,0x94,0x00,0x00,0x90,0x53,0xfb,0x0b,0x00,,...]
  }

```

Reply:
```
  {
    ID: <request ID>,
    replies: <array of UDP byte array>
  }

e.g.
  {
    ID: 19,
    replies: [
      [0x17,0x94,0x00,0x00,0x90,0x53,0xfb,0x0b,0xc0,0xa8,...],
      [0x17,0x94,0x00,0x00,0x41,0x78,0x1e,0x12,0xc0,0xa8,...],
    ]
  }
```

### HTTPS POST

The HTTPS POST connector is an HTTP POST connector that only accepts TLS client connections.

```
--in https/<bind address> [--html <folder>] [--ca-cert <file>] [--cert <file>] [--key <file>] [--client-auth]

  --html <folder> Folder containing the HTML served to the browser on the bind address.
  --ca-cert      CA certificate used to verify client certificates (defaults to ca.cert)
  --cert         server TLS certificate in PEM format (defaults to server.cert)
  --key          server TLS key in PEM format (defaults to server.key)
  --client-auth  requires client mutual authentication if supplied

e.g. 

--in https:/0.0.0.0:8080 --html examples/html
```

POST request:
```
  {
    ID: <request ID>,
    wait: <UDP timeout>,
    request: <UDP request byte array>
  }

e.g.

  {
    ID: 19,
    wait: "5s",
    request: [0x17,0x94,0x00,0x00,0x90,0x53,0xfb,0x0b,0x00,,...]
  }

```

Reply:
```
  {
    ID: <request ID>,
    replies: <array of UDP byte array>
  }

e.g.
  {
    ID: 19,
    replies: [
      [0x17,0x94,0x00,0x00,0x90,0x53,0xfb,0x0b,0xc0,0xa8,...],
      [0x17,0x94,0x00,0x00,0x41,0x78,0x1e,0x12,0xc0,0xa8,...],
    ]
  }
```


