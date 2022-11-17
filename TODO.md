# TODO

## IN PROGRESS

- [x] https://github.com/uhppoted/uhppoted-tunnel/issues/2
      - [x] Load TOML file
      - [x] Specify TOML file section on command line
      - [x] Use # or :: as seperator
      - [x] Override with command-line args
      - [x] Somehow roll lib.CommandX into lib.Command
      - [x] Clean up initial command configuration
      - [x] Merge back into _master_
      - [ ] Default to uhppoted-tunnel.toml
      - [ ] `daemonize`
            - [x] Use `config` for service args
            - [ ] Add label to TOML configuration
            - [x] Darwin
            - [ ] Linux
            - [ ] Windows
      - [ ] Update README
      - [ ] CONFIGURATION.md

- [ ] https://github.com/uhppoted/uhppoted-tunnel/issues/7
      - [x] Move lockfile implementation to uhppoted-lib
      - [x] Default to ephemeral _tmp_ folder for lockfiles
      - [ ] Use `flock` for default Linux implementation 

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

