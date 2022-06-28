# TODO

- [ ] README

- [x] HTTPS
      - [x] Update daemonize

- [ ] HTTP POST
      - [ ] Zero timeout for single requests
            - (but with internal timeout)
      - [ ] Use srv.Shutdown rather than Close (and remove related logic)

- [ ] Check tunnel -> tunnel -> tunnel interop
- [ ] Use global context for UDP broadcast timeout
- [ ] Make default max retries to be something other than 0
- [ ] More shutdown cleanup
```
2022/06/27 10:14:40 INFO   UDP    closing
2022/06/27 10:14:40 INFO   TCP    retrying in 5s
2022/06/27 10:14:40 INFO   TCP    closed
2022/06/27 10:14:40 INFO   UDP    retrying in 5s
2022/06/27 10:14:40 INFO   UDP    closed
```

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

