# TODO

- [ ] daemonize
      - [ ] daemonize client/server separately
      - [ ] portal/pipe

- [ ] Close()
      - [ ] shutdown cleanly with timeout
            - [x] router
            - [x] UDP broadcast
            - [x] UDP listen
            - [ ] TCP client
            - [ ] TCP server

- (?) Routing matrix
- (?) Replace handler functions with channels
- [ ] Remove dependency on uhppoted-lib and uhppote-core

## Miscellaneous

- [ ] HTTP POST portal
- [ ] gRPC portal
- [ ] Encode packet with protocol buffers
- [ ] undaemonize
- [ ] SSL + mutual auth
- [ ] Check interop with events
- (?) Wrap [libevent](https://libevent.org)
