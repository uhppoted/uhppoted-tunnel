# TODO

- [x] shutdown cleanly with timeout
      - [x] Fix closing/error/retry loop logic in UDP/TCP listen
            - [x] backoff
            - [x] error handling

- [ ] HTTP POST
      - [ ] close
      - [ ] ctx.cancel
      - [ ] request timeout
      - [ ] UDP timeout
      - [ ] decode result
      - [ ] array of replies
      - [ ] side-by-side debug

- [ ] Check tunnel -> tunnel -> tunnel interop

## Miscellaneous

- [ ] 'events' connectors
- [ ] gRPC portal
- [ ] WebUI

- (?) eBPF
- (?) Encode packet with protocol buffers
- (?) Wrap [libevent](https://libevent.org) or use syscalls
- (?) Routing matrix
- (?) Replace handler functions with channels
- (?) Remove dependency on uhppoted-lib and uhppote-core

