# `uhppoted-tunnel.toml`

`uhppoted-tunnel.toml` is a [TOML](https://toml.io/en/) file that defines the default configuration for _uhppoted-tunnel_.

The file is intended to allow simple and maintainable configuration for systems running multiple _uhppoted-tunnel_
services and provides a base configuration which can be overridden by command line arguments when necessary e.g. for 
running a _tunnel_ in console mode.

A sample `uhppoted-tunnel.toml` is included in the [examples](https://github.com/uhppoted/uhppoted-tunnel/blob/master/examples/uhppoted-tunnel.toml).

A `uhppoted-tunnel.toml` file comprises:
- a _defaults_ section which defines the base configuration for all _uhppoted-tunnels_
- service specific sections that customise the base configuration for a particular _uhppoted-tunnel_ (typically
  by defining the _in_ and _out_ connectors)

Running an instance of _uhppoted-tunnel_ with the `--config` flag selects the service specific section to use, e.g.
```
uhppoted-tunnel --config "#client" ...
```
will run with the [client] section from the default _uhppoted-tunnel.toml_ file.

```
uhppoted-tunnel --config "/etc/uhppoted/tunnel-special.toml#client" ...
```
will run with the [client] section from the _/etc/uhppoted/tunnel-special.toml#client_ file.

### Notes

1. The `daemonize` command supports the `--config` argument on the command line and will configure the _uhppoted-tunnel_
   daemon/service to run with the specified TOML file/section, e.g.
```
sudo ./uhppoted-tunnel daemonize --config "#client"
```

2. Command line arguments override any settings in the the TOML file, allowing _uhppoted-tunnel_ to be 
   easily run with custom configurations for e.g. testing and debugging:
```
./uhppoted-tunnel --config "#client" --debug --console --out udp/broadcast:192.168.1.255:60005 --udp-timeout 15s
```

3. _uhppoted-tunnel_ is configured on start only - changes to a TOML file will not take effect until the service/instance
   is restarted.

## [defaults] section

The _[defaults]_ section in the TOML sets the base configuration for instances of _uhppoted-tunnel_. The settings defined
in the _[defaults]_ section can be individually overriden in service specific sections, e.g.
```
[defaults]
lockfile = ""
max-retries = 32
max-retry-delay = "5m"
udp-timeout = "5s"
interfaces = { "in" = "en3", "out" = "" }
rate-limit = 1
rate-limit-burst = 120
...


```

### Settings

| *Attribute*      | *Description*                                                   | *Default value*                   |
| -----------------| ----------------------------------------------------------------|-----------------------------------|
| in               | _IN_ connector that accepts external requests                   | _None_                            |
| out              | _OUT_ connector that dispatches received requests               | _None_                            |
| interfaces       | _IN_ and _OUT_ connector network interfaces                     | _None_                            |
| lockfile         | lockfile used to prevent running multiple copies of the service | _auto-generated_                  |
| max-retries      | Maximum number of times to retry failed connection.             | -1 (retry forever)                |
| max-retry-delay  | Maximum delay between retrying failed connections               | 5m                                |
| udp-timeout      | Maximum delay between retrying failed connections               | 5s                                |
| ca-cert          | (TLS only) File path for CA certificate PEM file                | ./ca.cert                         |
| cert             | (TLS only) File path for client/server certificate PEM file     | ./client.cert or ./server.cert    |
| key              | (TLS only) File path for client/server key PEM file             | ./client.key  or ./server.key     |
| client-auth      | (TLS only) Mandates client authentication                       | false                             |
| authorisation    | (Tailscale only) Tailscale authorisation method                 | _TS_AUTHKEY_ environment variable |
| html             | (HTTP only) Folder with HTML                                    | ./html                            |
| log-level        | Sets the logging level (debug, info, warn or error)             | info./html                        |
| console          | Runs in _console_ mode i.e. logs to console                     | false                             |
| debug            | Enables display of low-level UDP messages                       | false                             |
| label            | Service label used to distinguish multiple tunnesl on a machine | _None_                            |
|                  |                                                                 |                                   |
| rate-limit       | Average request rate limit (requests/second)                    | 1                                 |
| rate-limit-burst | Burst request rate limit (requests)                             | 120                               |


## Service specific sections

A _service specific section_ defines the custom configuration (typically at least the _IN_ and _OUT_ connectors) for a 
particular service/instance of _uhppoted-tunnel_, e.g.:
```
...

[client]
in = "tcp/client:127.0.0.1:12345"
out = "udp/broadcast:192.168.1.255:60000"
udp-timeout = "1s"
log-level = "warn"
label = "qwerty"
...

```

The section name is specified in the `--config` command line argument
preceded by a _#_, e.g.
```
./uhppoted-tunnel --config "#client" 
```

## Tailscale authorisation

By default connections to a Tailscale tailnet will use the authorisation key in the TS_AUTHKEY environment variable. If the 
environment variable is not defined or is blank then you will be prompted with an authorisation URL. Alternative authorisation
methods can be configured using the `authorisation` TOML configuration file key:

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

3. Authorisation can be done using an OAuth2 client ID created on the the Tailscale admin console.

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
- The tailnet is the user or organisation account name ([**not** the tailnet DNS 
  name](https://github.com/tailscale/terraform-provider-tailscale/issues/206)) but can be defaulted to a '-' since the API
  keys are organisation/client specific.

Please note that connections authorised using _OAuth2_ are required to be _tagged_ and the keys do not expire, but 
can be expired manually on the Tailscale console.

## Sample TOML file

```
[defaults]
lockfile = ""
max-retries = 32
max-retry-delay = "5m"
udp-timeout = "5s"
interfaces = { "in" = "en3", "out" = "" }

[host]
in = "udp/listen:0.0.0.0:60000"
out = "tcp/server:0.0.0.0:12345"
lockfile = "/tmp/uhppoted-tunnel-host.pid"
label = "uiop"

[client]
in = "tcp/client::lo0:127.0.0.1:12345"
out = "udp/broadcast:192.168.1.255:60000"
udp-timeout = "1s"
log-level = "info"
label = "qwerty"
...
...
```