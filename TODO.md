# TODO

## IN PROGRESS

- [x] https://github.com/uhppoted/uhppoted-tunnel/issues/3
      - [x] TCP
        - [x] linux
        - [x] MacOS
        - [x] ~~Windows~~
      - [x] TLS
        - [x] linux
        - [x] MacOS
        - [x] ~~Windows~~
      - [x] UDP
        - [x] linux
        - [x] MacOS
        - [x] ~~Windows~~
      - [x] Add to TOML configuration
      - [x] README
      - [x] Help/usage

- [ ] TLS client to non-TLS host handshake doesn't ever timeout/disconnect

- [ ] 'events' connectors
      (?) let ID 0 imply no reply expected

## TODO

- io_uring
  - https://unixism.net/loti/index.html
- Socket activation
   - https://blog.podman.io/2023/01/systemd-socket-activation-lesson-learned/
- (?) https://eli.thegreenplace.net/2022/ssh-port-forwarding-with-go/
- (?) [UDP tunnelling: ssh/nc](https://superuser.com/questions/53103/udp-traffic-through-ssh-tunnel)
- (?) [UDP tunnelling: socat](http://www.morch.com/2011/07/05/forwarding-snmp-ports-over-ssh-using-socat/)

- [ ] Consider using IP_FREEBIND sockopt
      - https://www.freedesktop.org/wiki/Software/systemd/NetworkTarget

- [ ] https://tls-anvil.com/docs/Quick-Start/index

- [ ] Commonalize connector behaviours
- [ ] Use cancelable contexts throughout
- [ ] gRPC
- [ ] WSS
- [ ] XMPP
- [ ] ZMQ

- (?) eBPF
- (?) Encode packet with protocol buffers/MessagePack/Apache Avro
- (?) Wrap [libevent](https://libevent.org) or use syscalls
- (?) Routing matrix
- (?) Replace handler functions with channels
- (?) Remove dependency on uhppoted-lib and uhppote-core
- [ ] httpd connector logo
      - https://graphicdesign.stackexchange.com/questions/159149/how-to-draw-parallel-inclined-surfaces-in-perspective

