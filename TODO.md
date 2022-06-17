# TODO

- [ ] SSL + mutual auth

- [ ] shutdown cleanly with timeout
      - [x] router
      - [x] UDP broadcast
      - [x] UDP listen
      - [x] TCP server
      - [x] TCP client
            - [x] Close is slow when still trying to connect i.e. host not running
      - [ ] Ignore warnings if closing
      - [ ] Fix closing/error/retry loop logic in UDP/TCP listen
            - [ ] backoff
            - [ ] error handling
            - (?) Condition handler a la LISP

- [ ] Remove dependency on uhppoted-lib and uhppote-core

## Miscellaneous

- [ ] HTTP POST portal
- [ ] gRPC portal
- [ ] Check interop with events
- (?) Encode packet with protocol buffers
- (?) Wrap [libevent](https://libevent.org) or use syscalls
- (?) Routing matrix
- (?) Replace handler functions with channels
