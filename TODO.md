# TODO

- [ ] README
- [x] Move `html` folder to `examples` or somesuch

- [ ] Check tunnel -> tunnel -> tunnel interop
- [ ] Use global context for UDP broadcast timeout
- [x] Make default max retries to be something other than 0
- [x] More shutdown cleanup
- [ ] Remove lockfile on fatalf (e.g. retries)
      - (?) shutdown hook maybe

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

