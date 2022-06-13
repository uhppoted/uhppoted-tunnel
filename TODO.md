# TODO

- [ ] daemonize
      - [ ] portal/pipe
      - [x] daemonize client/server separately
      - [ ] UFW firewall rules

- [ ] undaemonize
    - [x] MacOS
    - [ ] Linux
    - [ ] Windows

- [ ] Close()
      - [ ] shutdown cleanly with timeout
            - [x] router
            - [x] UDP broadcast
            - [x] UDP listen
            - [ ] TCP client
            - [ ] TCP server

- [ ] Remove dependency on uhppoted-lib and uhppote-core

## Miscellaneous

- [ ] SSL + mutual auth
- [ ] HTTP POST portal
- [ ] gRPC portal
- [ ] Check interop with events
- (?) Encode packet with protocol buffers
- (?) Wrap [libevent](https://libevent.org)
- (?) Routing matrix
- (?) Replace handler functions with channels
