# CHANGELOG

## [0.8.9](https://github.com/uhppoted/uhppoted-tunnel/releases/tag/v0.8.9) - 2024-09-06

### Added
1. _ip/out_ connector that supports UDP broadcast, UDP direct connections and TCP connections
   to controllers.

### Updated
1. Updated to Go 1.23.


## [0.8.8](https://github.com/uhppoted/uhppoted-tunnel/releases/tag/v0.8.8) - 2024-03-27

### Added
1. Added `restore-default-parameters` command to the HTTP connector example.

### Updated
1. Bumped Go version to 1.22


## [0.8.7](https://github.com/uhppoted/uhppoted-tunnel/releases/tag/v0.8.7) - 2023-12-01

### Added
1. Added `set-door-passcodes` command to the HTTP connector example.


## [0.8.6](https://github.com/uhppoted/uhppoted-tunnel/releases/tag/v0.8.6) - 2023-08-30

### Added
1. Added `activate-keypads` command to the HTTP connector example.


## [0.8.5](https://github.com/uhppoted/uhppoted-tunnel/releases/tag/v0.8.5) - 2023-06-13

### Added
1. Implemented tailscale client and server connectors.
2. Added rate limiter for incoming requests/events.
3. Added `set-interlock` command to the HTTP connector example

### Updated
1. Renamed _master_ branch to _main_.


## [0.8.4](https://github.com/uhppoted/uhppoted-tunnel/releases/tag/v0.8.4) - 2023-03-17

### Added
1. Support for binding to a specific network interface (Linux and MacOS only). 
   Ref. [Pi Configuration as secure wireless tunnel](https://github.com/uhppoted/uhppoted-tunnel/issues/3)
2. UDP, TCP and TLS _event_ connectors for relaying events.
3. `doc.go` package overview documentation.
4. Added card PIN field to HTTP connector example.

### Updated
1. Fixed timeout in TCP and TLS clients
2. Fixed 'fatal' error on closed


## [0.8.3](https://github.com/uhppoted/uhppoted-tunnel/releases/tag/v0.8.3) - 2022-12-16

### Added
1. TOML file configuration (cf. https://github.com/uhppoted/uhppoted-tunnel/issues/2)
2. Experimental Linux/ARM64 binary

### Changed
1. Moved default lockfile to ephemeral system _temp_ folder
2. Updated Linux _systemd_ unit file to wait on `network-online.target` (ref. https://systemd.io/NETWORK_ONLINE)


## [0.8.2](https://github.com/uhppoted/uhppoted-tunnel/releases/tag/v0.8.2) - 2022-10-14

### Added
1. HTTP: implemented remaining UHPPOTE functions
2. HTTP: reworked to used codegen'd Javascript
3. Commonalised HTTP/S implementations

### Changed
1. Changed _send_ URL to `/udp/send`
2. Added _broadcast_ URL `/udp/broadcast`
3. Fixed typo in service name.
4. Fixed formatting in daemonize command.
5. Fixed missing mutex in router.go and restructured (cf. https://github.com/uhppoted/uhppoted-tunnel/issues/4)


## [0.8.1](https://github.com/uhppoted/uhppoted-tunnel/releases/tag/v0.8.1) - 2022-08-01

### Changed
1. Added `set-address` implementation to HTTP/S example console
2. Added `get-time` implementation to HTTP/S example console
3. Added `set-time` implementation to HTTP/S example console


## [0.8.0](https://github.com/uhppoted/uhppoted-tunnel/releases/tag/v0.8.0) - 2022-07-01

1. Initial release

