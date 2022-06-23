# TODO

- [ ] shutdown cleanly with timeout
      - [x] Ignore warnings if closing
      - [x] Fix delay in TCP/TLS server retry
      - [ ] Fix closing/error/retry loop logic in UDP/TCP listen
            - [ ] backoff
            - [ ] error handling
            - (?) Condition handler a la LISP

- [ ] HTTP POST
- [ ] Check interop with events
- [ ] Check tunnel -> tunnel -> tunnel interop

## Miscellaneous

- (?) eBPF
- (?) gRPC portal
- (?) Encode packet with protocol buffers
- (?) Wrap [libevent](https://libevent.org) or use syscalls
- (?) Routing matrix
- (?) Replace handler functions with channels
- (?) Remove dependency on uhppoted-lib and uhppote-core

