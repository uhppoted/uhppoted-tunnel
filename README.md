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
- UDP event
- TCP server
- TCP client
- TLS server
- TLS client
- HTTP POST
- HTTPS POST
- Tailscale server
- Tailscale client
- IP/out

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
- Linux/ARM64 (experimental)

## Release Notes

### Current Release

**[v0.8.9](https://github.com/uhppoted/uhppoted-tunnel/releases/tag/v0.8.9) - 2024-09-06

1. Added _ip/out_ connector that supports UDP broadcast, UDP direct connections and TCP connections
   to controllers.
2. Updated to Go 1.23.


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
- [Go 1.20+](https://go.dev)
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

| *Dependency*                                                      | *Description*                            |
|-------------------------------------------------------------------|------------------------------------------|
| [uhppote-core](https://github.com/uhppoted/uhppote-core)          | Device level API implementation          |
| [uhppoted-lib](https://github.com/uhppoted/uhppoted-lib)          | Common library functions                 |
| golang.org/x/sys                                                  | (for Windows service integration)        |
| tailscale.com                                                     | _tsnet_ library for Tailscale connectors |

## uhppoted-tunnel

Usage: ```uhppoted-tunnel <command> --in <connector> --out <connector> <options>```

Supported commands:

- `help`
- `version`
- `run`
- `daemonize`
- `undaemonize`

Defaults to `run` if the command it not provided i.e. ```uhppoted-tunnel --in <connector> --out <connector> <options>``` is equivalent to ```uhppoted-tunnel run  --in <connector> --out <connector> <options>```.

#### Configuration

For _uhppoted-tunnel_ v0.8.3+, runtime configuration is defined in a TOML file (documented [here](https://github.com/uhppoted/uhppoted-tunnel/blob/master/documentation/uhppoted-tunnel-toml.md)) and any future enhancements will
be configurable only in the TOML file.

The command line arguments described below are for legacy support and overriding specific settings in the TOML configuation.

### `run`

Runs the `uhppoted-tunnel` service. Default command, intended for use as a system service that runs in the 
background. 

Command line:

` uhppoted-tunnel [--debug] [--console] --config <configuration> --in <connector> --out <connector> [options]`

```
  --config <configuration> Sets the TOML file and section to use for runtime configuration settings. The
                           configuration may be:
                           - fully specified, e.g. "--config /etc/uhppoted/uhppoted-tunnel.toml#client"
                           - file only e.g. "--config /etc/uhppoted/uhppoted-tunnel.toml" (uses the [defaults] section)
                           - section only e.g. "--config #client" (uses the default TOML file and [client] section)
                           
                           If the --config argument not supplied, the default TOML file will be used if it exists.

  --in <connector>  Defines the connector that accepts incoming commands. Overrides the 'IN' connector in the TOML
                    configuration if it exists. Valid 'in' connectors include: 
                    - udp/listen:<bind address> (e.g. udp/listen:0.0.0.0:60000)
                    - udp/event:<bind address> (e.g. udp/listen:0.0.0.0:60000)
                    - tcp/server:<bind address> (e.g. tcp/server:0.0.0.0:12345)
                    - tcp/client:<host address> (e.g. tcp/client:192.168.1.100:12345)
                    - tls/server:<bind address> (e.g. tls/server:0.0.0.0:12345)
                    - tls/client:<host address> (e.g. tls/client:192.168.1.100:12345)
                    - tailscale/server:<server address> (e.g.uhppoted:12345,nolog)
                    - http/<bind address> (e.g. http/0.0.0.0:8080)
                    - https/<bind address> (e.g. https/0.0.0.0:8443)

                    Under Linux and MacOS TCP and UDP _in_ connectors can be bound to a specific interface by prefixing
                    the address with ::<interface> e.g. tcp/client::en3:192.168.1.100:12345. The _Tailscale_ connector
                    syntax is described below.

  --out <connector> Defines the connector that forwards received commands. Overrides the 'OUT' connector in the TOML
                    configuration if it exists. Valid 'out' connectors include: 
                    - udp/broadcast:<broadcast address> (e.g. udp/broadcast:255.255.255.255:60000)
                    - udp/event:<broadcast address> (e.g. udp/broadcast:255.255.255.255:60000)
                    - tcp/server:<bind address> (e.g. tcp/server:0.0.0.0:12345)
                    - tcp/client:<host address> (e.g. tcp/client:192.168.1.100:12345)
                    - tls/server:<bind address> (e.g. tls/server:0.0.0.0:12345)
                    - tls/client:<host address> (e.g. tls/client:192.168.1.100:12345)
                    - tailscale/client:<client address> (e.g. tailscale/client::makerspace:uhppoted:12345,nolog)

                    Under Linux and MacOS TCP and UDP _out_ connectors can be bound to a specific interface by prefixing
                    the address with ::<interface> e.g. udp/broadcast::lo0:127.0.0.01:12345. The _Tailscale_ connector
                    syntax is described below.

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

  --ca-cert <file>  (TLS only) File path for CA certificate PEM file. Defaults to ./ca.cert

  --cert <file>     (TLS only) File path for client/server certificate PEM file. Defaults to./client.cert ('IN' 
                               connectors) or ./server.cert (OUT connectors)

  --key <file>      (TLS only) File path for client/server key PEM file. Defaults to ./client.key ('IN' connectors)
                               or ./server.key ('OUT' connectors)
 
  --client-auth     (TLS only) Mandates client authentication. Defaults to false

  --html            (HTTP only) Folder with HTML, CSS, images, etc. Defaults to./html
```

In general, tunnels operate in pairs - one on the _host_, listening for commands from e.g. the _AccessControl_ application
or _uhppote-cli_ and the other on the _client_ local to the controller, which sends the commands to the controller(s)
and returns the replies to the _host_. It is however, possible to chain multiple tunnels to bridge across several machines.

The _event_ connectors are a specialization of the UDP _listen_ and _broadcast_ connectors in that events are relayed in a
a single direction only, without expectation of a reply. It's quite possible to use the _listen_ and _broadcast_ connectors
to relay events but the specialized connectors are slightly optimized for the use case and have also been put in place to 
support future enhancements that may rely on the specialized connectors.

### `daemonize`

Registers `uhppoted-tunnel` as a system service that will be started on system boot. The command creates the necessary
system specific service configuration files and service manager entries. 

On Linux:
- The service defaults to using the `uhppoted:uhppoted` user:group - this can be changed with the `--user` option
- Depending on the system, it may be necessary to run `sudo systemctl enable uhppoted-tunnel-xxx` after _daemonizing_
  to get the _uhppoted-tunnel_ service to start on boot.
- By default, the service is configured to wait for the `network-online.target` (cf. https://systemd.io/NETWORK_ONLINE). To wait
  for a specific interface modify the unit file (_/etc/systemd/system/uhpppoted-tunnel-xxx_) to wait for [systemd-networkd-wait-online.service](https://www.freedesktop.org/software/systemd/man/systemd-networkd-wait-online.service.html)

Command line:

`uhppoted-tunnel daemonize --config <configuration> --in <connector> --out <connector> [--label <label>] [--user <user>]`

```
  --config <configuration> Sets the TOML file and section to use for runtime configuration settings. The
                           configuration may be:
                           - fully specified, e.g. "--config /etc/uhppoted/uhppoted-tunnel.toml#client"
                           - file only e.g. "--config /etc/uhppoted/uhppoted-tunnel.toml" (uses the [defaults] section)
                           - section only e.g. "--config #client" (uses the default TOML file and [client] section)
                           
                           If the --config argument not supplied, the default TOML file will be used if it exists.

  --in <connector>  Defines the connector that accepts incoming commands. Overrides the 'in' connector in the TOML
                    configuration. Valid 'in' connectors include: 
                    - udp/listen:<bind address> (e.g. udp/listen:0.0.0.0:60000)
                    - tcp/server:<bind address> (e.g. tcp/server:0.0.0.0:12345)
                    - tcp/client:<host address> (e.g. tcp/client:192.168.1.100:12345)
                    - tls/server:<bind address> (e.g. tls/server:0.0.0.0:12345)
                    - tls/client:<host address> (e.g. tls/client:192.168.1.100:12345)
                    - http/<bind address> (e.g. http/0.0.0.0:8080)
                    - https/<bind address> (e.g. https/0.0.0.0:8443)

  --out <connector> Defines the connector that forwards received commands. Overrides the 'out' connector in the TOML
                    configuration. Valid 'out' connectors include: 
                    - udp/broadcast:<broadcast address> (e.g. udp/broadcast:255.255.255.255:60000)
                    - tcp/server:<bind address> (e.g. tcp/server:0.0.0.0:12345)
                    - tcp/client:<host address> (e.g. tcp/client:192.168.1.100:12345)
                    - tls/server:<bind address> (e.g. tls/server:0.0.0.0:12345)
                    - tls/client:<host address> (e.g. tls/client:192.168.1.100:12345)

  --label <label>  Identifying label for the tunnel daemon/service, used to identify the tunnel in logs and when
                   uninstalling the daemon/service. Imperative if running multiple tunnel daemons on the same machine,
                   optional but recommended otherwise. Defaults to uhppoted-tunnel if not provided.

  --user <uid:group>  (Linux only) uid:group pair to use for service. Defaults to uhppoted:uhppoted.
```

### `undaemonize`

Unregisters `uhppoted-tunnel` as a system service, but does not delete any created log or configuration files. 

Command line:

`uhppoted-tunnel undaemonize [--label <label>]`


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
- Tailscale server

_OUT_ connectors:

- UDP broadcast
- TCP server
- TCP client
- TLS server
- TLS client
- Tailscale client
- IP

### UDP listen

Listens for incoming UDP packets on the _bind address_, effectively acting as a direct proxy for a remote controller.

```
--in udp/listen[::<interface>]:<bind address>

e.g. 

--in udp/listen:0.0.0.0:60000
--in udp/listen::en3:0.0.0.0:60000
```

### UDP broadcast

Sends a received packet out as a UDP message on the _broadcast address_ and forwards any replies to the original requester,
effectively acting as a proxy for a remote application.

```
--out udp/broadcast[::<interface>]:<broadcast address> [--udp-timeout <timeout>]

   The broadcast address is typically (but not necessarily) the UDP broadcast for the network adapter for the controllers'
   network segment. However it can be any valid IPv4 address:port combination to accomodate the requirements of the 
   installation.

   --udp-timeout <timeout>  Sets the maximum time to wait for replies to a broadcast message, in human readable format
                            e.g. 15s, 1250ms, etc. Defaults to 5 seconds if not provided.

e.g. 

--out udp/broadcast:255.255.255.255:60000 --udp-timeout 5s
--out udp/broadcast::en3:255.255.255.255:60000 --udp-timeout 5s
```

### TCP server

The TCP server connector accepts connections from one or more TCP clients and can act as both an _IN_ connector and an _OUT_ connector.
Incoming requests will be forwarded to all connected clients.

```
--in tcp/server[::<interface>]:<bind address>

e.g. 

--in tcp/server:0.0.0.0:12345
--in tcp/server::en3:0.0.0.0:12345
```

### TCP client

The TCP client connector connects to a TCP server and can act as both an _IN_ connector and an _OUT_ connector. Incoming requests/replies
will be forwarded to the remote server.

```
--in tcp/client[::<interface>]:<host address>

e.g. 

--in tcp/host:192.168.1.100:12345
--in tcp/host::lo0:127.0.0.1:12345
```

### TLS server

The TLS server connector is a TCP server connector that only accepts TLS secured client connections.

```
--in tls/server[::<interface>]:<bind address> [--ca-cert <file>] [--cert <file>] [--key <file>] [--client-auth]

  --ca-cert      CA certificate used to verify client certificates (defaults to ca.cert)
  --cert         server TLS certificate in PEM format (defaults to server.cert)
  --key          server TLS key in PEM format (defaults to server.key)
  --client-auth  requires client mutual authentication if supplied

e.g. 

--in tls/server:0.0.0.0:12345 --ca-cert tunnel.ca --cert tunnel.cert --key tunnel.key --client-auth
--in tls/server::en3:0.0.0.0:12345 --ca-cert tunnel.ca --cert tunnel.cert --key tunnel.key --client-auth
```

### TLS client

The TLS client connector is a TCP client connector that only connects to TLS secured servers.

```
--in tls/client[::<interface>]:<host address> [--ca-cert <file>] [--cert <file>] [--key <file>] [--client-auth]

  --ca-cert      CA certificate used to verify server certificates (defaults to ca.cert)
  --cert         client TLS certificate in PEM format. Optional, only required if the TLS server 
                 has mutual authentication enabled.
  --key          client TLS key in PEM format. Optional, only required if the TLS server 
                 has mutual authentication enabled.

e.g. 

--in tls/client:192.168.1.100:12345 --ca-cert tunnel.ca --cert client.cert --key client.key
--in tls/client::en3:192.168.1.100:12345 --ca-cert tunnel.ca --cert client.cert --key client.key
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

### _Tailscale_ 

#### _Tailscale_ server

The _Tailscale_ server connector establishes a listening Tailscale connection that accepts connections from one or more Tailscale
clients. Unlike the TCP connectors it is designed to act only as an _IN_ connector. Incoming requests are be forwarded to all connected clients.

The Tailscale server _address_ comprises:
- tailnet `machine` name to use for the server (required)
- IP port on which to accept incoming connections
- optional `debug` to display the tailscale connection logging. Defaults to 'no log' because the Tailscale logging is 
  very verbose, but occasionally useful or necessary for debugging connection issues. The only valid value is `debug` -
  other values (e.g. nolog) can be used as placeholder mnemonics.

```
--in tailscale/server:<device>:<port>[,debug]

e.g. 

--in tailscale/server:uhppoted:12345,debug
--in tailscale/server:uhppoted:12345,nolog
```

Please note that Tailscale does not currently suppport binding to a specific interface (Ref. https://github.com/tailscale/tailscale/issues/1552).

#### _Tailscale_ client

The _Tailscale_ client connector connects to a listening _Tailscale_ server connecton. Unlike the TCP and UDP client connectors,
a Tailscale client connector can only be configured as an _OUT_ connector i.e. it expects to receive commands on the UDP IN port
and forwards the commands to the Tailscale server that is connected to the access controller.

The Tailscale client _address_ comprises:
- tailnet `machine` name to use for the server (optional, but recommended)
- Tailnet server to which to connect, in the format <machine>:<port> (required)
- optional `debug` to display the tailscale connection logging. Defaults to 'no log' because the Tailscale logging is 
  very verbose, but occasionally useful or necessary for debugging connection issues. The only valid value is `debug` -
  other values (e.g. nolog) can be used as placeholder mnemonics.


```
--in tailscale/client[::<machine>]:<server address>[,debug]

e.g. 

--out tailscale/client:uhppoted:12345
--out tailscale/client::qwerty:uhppoted:12345,debug
--out tailscale/client::qwerty:uhppoted:12345,nolog
```

Please note that Tailscale does not currently suppport binding to a specific interface (Ref. https://github.com/tailscale/tailscale/issues/1552).


#### Authorisation

By default connections to a Tailscale tailnet will use the authorisation key in the TS_AUTHKEY environment variable. If the environment variable is not defined or is blank then you will be prompted with an authorisation URL.

Alternative authorisation methods can be configured in the TOML configuration file (specified using the `--config` command
line option):

1. A different environment variable can specified using the `env:<variable name>` syntax, e.g.
```
[tailscale-server]
...
authorisation = "env:TS_WORKSHOP"
...
```
   This is an alternative to using a reusable authorisation key in the TS_AUTHKEY environment variable when running two
   or more tunnels on the same machine.

2. The authorisation key can specified directly using the `authkey:<key>` syntax, e.g.
```
[tailscale-server]
...
authorisation = "authkey:tskey-auth-xxxxxxxxxxxx-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
...
```

3. An OAuth2 client, using the `oauth2:<credentials>` syntax, e.g.
```
[tailscale-server]
...
authorisation = "oauth2:.credentials.workship"
...
```

The `credentials` is a JSON file that contains the OAuth2 credentials for the OAuth2 client, e.g.
```
{ 
    "tailscale": {
        "oauth2": {
            "client-id": "xxxxxxxxxxxx",
            "client-secret": "tskey-client-xxxxxxxxxxxx-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
            "auth-url": "https://api.tailscale.com/api/v2/oauth/token",
            "tailnet": "qwerty@uiop.com",
            "tag": "development",
            "key-expiry": 300
        }
    }
}
```

- The `client-id` and `client-secret` are the keys generated when creating the OAuth2 client on the Tailscale admin console.
- The `tailnet` is the user/organisation account name ([**not** the tailnet DNS name](https://github.com/tailscale/terraform-provider-tailscale/issues/206)) but can be defaulted to a '-' since the API keys are organisation/client specific.
- Connections authorised using _OAuth2_ are required to be _tagged_ and the keys do not expire (but can be expired manually
  on the Tailscale console.)

#### Notes

1. _tailscale/client_ connectors are declared as ephemeral i.e. will be cleaned up when disconnected.
2. _tailscale/server_ connectors are not ephemeral. This is necessary for the case where a server restarts and
   needs to reconnect as the 'same' machine so that existing clients can reconnect without having to be 
   restarted.
3. There is no internal 'keep-alive' - if a _tailscale-server_ goes offline the _tailscale-client_ will only
   become aware of it when forwarding a command received on the UDP IN connector. This is by design - an application
   specific keep-alive is more flexible and more useful. Additionally an internal keep-alive can potentially force
   a connection to unnecessarily backoff to the maximum delay.

### IP/out

The IP/out connector supports the following connections to controllers:
- UDP broadcast
- UDP 'sendto' connections
- TCP connections.

The default connection is UDP broadcast but specific controllers in the TOML configuration file can be configured for
UDP 'sendto' or TCP connections, making it a 'TOML' only connector for all practical purposes:

```
uhppoted-tunnel --config "uhppoted-tunnel.toml#ip"

where the TOML 'ip' section comprises:
...
[ip]
in = "udp/listen:0.0.0.0:60000"
out = "ip/out:192.168.1.255:60005"
console=true
debug = true

    [ip.controllers]
    405419896 = "udp::192.168.1.100:60005"
    303986753 = "tcp::192.168.1.100:60005"
...

- the 'in' connection is any supported IN connection
- the 'out' connection defines the default UDP broadcast connection
- the [controllers] subsection lists the controllers with transport protocol and IPv4 address
```


### _Rate Limiting_ 

_uhppoted-tunnel_ has an internal rate limit that limits the number of requests per second that can be processed. The default
settings are reasonably conservative:
- rate limited to 1 request per second
- burst rate of 120 requests

These limits can be changed in the TOML configuration file, e.g.:
```
...
rate-limit = 5
rate-limit-burst = 300
...
```

Fractional rate limits are supported e.g. `rate-limit = 0.1`

## Attribution

1. HTTP/S connector example logo uses [image](https://www.freepik.com/free-photo/light-shine-through-round-holes-ceiling-casting-shadows_15317209.htm) 
designed by [Garry Killian](https://www.freepik.com/author/garrykillian) for [freepik.com](https://www.freepik.com).
