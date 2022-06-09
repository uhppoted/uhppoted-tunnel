# TODO

- [x] lockfiles
- [x] Make UDP broadcast timeout configurable
- (?) Log to logfile
- [x] --log-level

- [ ] Error
```
2022/06/09 15:10:53 WARN   TCP    msg 26  error sending message to 127.0.0.1:50371 (write tcp 127.0.0.1:12345->127.0.0.1:50371: write: 
broken pipe)
```

- [ ] Close()
      - [ ] shutdown cleanly with timeout
            - [x] router
            - [ ] UDP broadcast
            - [ ] UDP listen
            - [ ] TCP client
            - [ ] TCP server

- [ ] Routing matrix
- (?) Replace handler functions with channels

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
