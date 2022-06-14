# TODO

- [ ] daemonize
      - [x] Make --label optional but recommended
      - [ ] Linux
            - [x] UFW firewall rules
      - [ ] Windows

- [ ] undaemonize
    - [x] MacOS
    - [ ] Linux
    - [ ] Windows

- [ ] Close()
      - [ ] shutdown cleanly with timeout
            - [x] router
            - [x] UDP broadcast
            - [x] UDP listen
            - [ ] TCP client
            - [x] TCP server
            - [ ] Fix closing/error/retry loop logic in UDP/TCP listen
                  - [ ] backoff
                  - [ ] error handling
                  - (?) Condition handler a la LISP

- [ ] Remove dependency on uhppoted-lib and uhppote-core

## Miscellaneous

- [ ] SSL + mutual auth
- [ ] HTTP POST portal
- [ ] gRPC portal
- [ ] Check interop with events
- (?) Encode packet with protocol buffers
- (?) Wrap [libevent](https://libevent.org)
- (?) Routing matrix
- (?) Replace handler functions with channels
