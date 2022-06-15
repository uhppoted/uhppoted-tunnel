# TODO

- [ ] daemonize
      - [x] --label optional but recommended
            - (?) lockfile
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
            - [x] TCP client
            - [x] TCP server
            - [ ] Ignore warnings if closing
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
