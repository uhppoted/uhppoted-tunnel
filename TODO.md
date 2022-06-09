# TODO

- [x] lockfiles

- [ ] Close()
      - [ ] shutdown cleanly with timeout
            - [x] router
            - [ ] UDP broadcast
            - [ ] UDP listen
            - [ ] TCP client
            - [ ] TCP server

- [ ] Routing matrix
- (?) Replace handler functions with channels
- [ ] Make UDP broadcast timeout configurable

## Miscellaneous
- [ ] HTTP POST portal
- [ ] gRPC portal
- [ ] Remove dependency on uhppoted-lib and uhppote-core
- [ ] Encode packet with protocol buffers
- [ ] daemonize
- [ ] undaemonize
- [ ] SSL + mutual auth
- [ ] Check interop with events
- (?) Wrap [libevent](https://libevent.org)
