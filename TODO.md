# TODO

- [x] README
- [x] Check tunnel -> tunnel -> tunnel interop
- [x] Remove lockfile on fatalf (e.g. after max retries)

## Miscellaneous

- [ ] Commonalize connector behaviours
- [ ] Use cancelable contexts throughout
- [ ] 'events' connectors
      (?) let ID 0 imply no reply expected
- [ ] gRPC portal
- [ ] WebUI

- (?) eBPF
- (?) Encode packet with protocol buffers
- (?) Wrap [libevent](https://libevent.org) or use syscalls
- (?) Routing matrix
- (?) Replace handler functions with channels
- (?) Remove dependency on uhppoted-lib and uhppote-core

