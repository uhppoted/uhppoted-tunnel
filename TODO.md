# TODO

## IN PROGRESS

- [x] https://github.com/uhppoted/uhppoted-tunnel/issues/2
      - [ ] `daemonize`
            - [x] `daemonize --in  udp/listen:0.0.0.0:60000 ....` without default uhppoted-tunnel.conf
            - [x] Use `config` for service args
            - [x] Add label to TOML configuration
            - [x] Darwin
            - [x] Linux
            - [x] Windows
            - [x] Add [console-client] etc to example TOML file (so that services don't have --console enabled by mistake)
            - [x] Commonalise the 'no label' handling
            - [x] Remove [console-client] etc and replace with --console in the Makefile
            - [ ] Commonalise config file resolution
            - [x] Enable --service for Linux and MacOS
      - [x] uhppoted-tunnel-toml.md
      - [x] Update README

- [ ] https://github.com/uhppoted/uhppoted-tunnel/issues/7
      - [x] Move lockfile implementation to uhppoted-lib
      - [x] Default to ephemeral _tmp_ folder for lockfiles
      - [ ] Use `flock` for default Linux implementation 
      - [ ] (optional) soft lock
      - [ ] (optional) socket lock

- [x] ARM64 build
- (?) https://eli.thegreenplace.net/2022/ssh-port-forwarding-with-go/
- [ ] log.Warnf+ should default to stderr
- [ ] Windows eventlog message file
      - https://social.msdn.microsoft.com/Forums/windowsdesktop/en-US/deaa0055-7770-4e55-a5b8-6d08b80b74af/creating-event-log-message-files

## TODO

- (?) [UDP tunnelling: ssh/nc](https://superuser.com/questions/53103/udp-traffic-through-ssh-tunnel)
- (?) [UDP tunnelling: socat](http://www.morch.com/2011/07/05/forwarding-snmp-ports-over-ssh-using-socat/)

- [ ] 'events' connectors
      (?) let ID 0 imply no reply expected
- [ ] https://tls-anvil.com/docs/Quick-Start/index

- [ ] Commonalize connector behaviours
- [ ] Use cancelable contexts throughout
- [ ] gRPC
- [ ] WSS
- [ ] XMPP

- (?) eBPF
- (?) Encode packet with protocol buffers
- (?) Wrap [libevent](https://libevent.org) or use syscalls
- (?) Routing matrix
- (?) Replace handler functions with channels
- (?) Remove dependency on uhppoted-lib and uhppote-core
- [ ] httpd connector logo
      - https://graphicdesign.stackexchange.com/questions/159149/how-to-draw-parallel-inclined-surfaces-in-perspective

