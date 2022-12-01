# CHANGELOG

## [Unreleased]

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

