# TODO

- [ ] README

- [ ] HTTP POST
      - [x] FS from command line
      - [x] array of replies
      - [x] UDP timeout
      - [x] request timeout
      - [x] listen: retry with backoff
      - [ ] Use global context from httpd struct
      - [ ] side-by-side debug
      - [x] eslint

- [ ] HTTPS

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

