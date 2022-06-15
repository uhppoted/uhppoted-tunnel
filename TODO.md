# TODO

- [ ] daemonize
      - [x] --label optional but recommended
            - (?) lockfile
      - [x] Linux
            - [x] UFW firewall rules
      - [ ] Windows

- [ ] undaemonize
    - [x] MacOS
    - [x] Linux
          - [x] `ERROR: remove /etc/logrotate.d/uhppoted-tunnnel-uiop: no such file or directory`
          - [x] `/var/uhppoted/tunnel/var/uhppoted/uhppoted-tunnnel-qwerty.pid.pid`
    - [ ] Windows

- [ ] Close()
      - [ ] shutdown cleanly with timeout
            - [x] router
            - [x] UDP broadcast
            - [x] UDP listen
            - [x] TCP server
            - [ ] TCP client
                  - [ ] Close is slow when still trying to connect i.e. host not running
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
