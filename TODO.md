# TODO

## IN PROGRESS

- [ ] https://github.com/uhppoted/uhppoted-tunnel/issues/3
      - [x] TCP
        - [x] linux
        - [x] MacOS
        - [x] ~~Windows~~

      - TLS
        - [ ] linux
        - [ ] MacOS
        - [x] ~~Windows~~

      - UDP
        - [ ] linux
        - [ ] MacOS
        - [x] ~~Windows~~

## TODO

- (?) https://eli.thegreenplace.net/2022/ssh-port-forwarding-with-go/
- (?) [UDP tunnelling: ssh/nc](https://superuser.com/questions/53103/udp-traffic-through-ssh-tunnel)
- (?) [UDP tunnelling: socat](http://www.morch.com/2011/07/05/forwarding-snmp-ports-over-ssh-using-socat/)

- [ ] Consider using IP_FREEBIND sockopt
      - https://www.freedesktop.org/wiki/Software/systemd/NetworkTarget

- [ ] 'events' connectors
      (?) let ID 0 imply no reply expected
- [ ] https://tls-anvil.com/docs/Quick-Start/index

- [ ] Commonalize connector behaviours
- [ ] Use cancelable contexts throughout
- [ ] gRPC
- [ ] WSS
- [ ] XMPP

- (?) eBPF
- (?) Encode packet with protocol buffers
- (?) Wrap [libevent](https://libevent.org) or use syscalls
- (?) Routing matrix
- (?) Replace handler functions with channels
- (?) Remove dependency on uhppoted-lib and uhppote-core
- [ ] httpd connector logo
      - https://graphicdesign.stackexchange.com/questions/159149/how-to-draw-parallel-inclined-surfaces-in-perspective

