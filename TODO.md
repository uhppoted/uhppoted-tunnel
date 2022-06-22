# TODO

- [x] Reorganize as --in and --out
- [x] Generalise backoff/retry
- [ ] Moved packet ID stuff to Conn

- [ ] shutdown cleanly with timeout
      - [ ] Ignore warnings if closing
      - [ ] Fix delay in TCP/TLS server retry
      - [ ] Fix closing/error/retry loop logic in UDP/TCP listen
            - [ ] backoff
            - [ ] error handling
            - (?) Condition handler a la LISP

## Miscellaneous

- [ ] HTTP POST portal
- [ ] gRPC portal
- [ ] Check interop with events
- (?) Encode packet with protocol buffers
- (?) Wrap [libevent](https://libevent.org) or use syscalls
- (?) Routing matrix
- (?) Replace handler functions with channels
- (?) Remove dependency on uhppoted-lib and uhppote-core

