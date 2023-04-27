# TODO

## IN PROGRESS

- [x] Rate limit (for e.g. when you've misconfigured your UDP broadcast to send to your own listener)
      - [x] Rate limiter
      - [x] TOML file settings
      - [x] README
      - [x] CHANGELOG

- [x] Error on CTRL-C
- [ ] Document openssl commands to generate keys
      - Ref. https://openssl-ca.readthedocs.io/en/latest/sign-server-and-client-certificates.html
      - Ref. https://stackoverflow.com/questions/10175812/how-to-generate-a-self-signed-ssl-certificate-using-openssl

- [ ] Tailscale
      - [x] Close server more gracefully
      ```
      2023/04/26 11:06:19 WARN   tailscale  read tcp 100.112.82.254:12345: operation aborted
      ```
      - [x] socket is closed in listen() - probaby don't need to close in sockets.Close as well
      - [ ] Update TCP/TLS server to match
      - [ ] plugin (in branch)


## TODO

- (?) https://shadowsocks.org

- [ ] integration tests
- io_uring
  - https://unixism.net/loti/index.html
- Socket activation
   - https://blog.podman.io/2023/01/systemd-socket-activation-lesson-learned/
- (?) https://eli.thegreenplace.net/2022/ssh-port-forwarding-with-go/
- (?) [UDP tunnelling: ssh/nc](https://superuser.com/questions/53103/udp-traffic-through-ssh-tunnel)
- (?) [UDP tunnelling: socat](http://www.morch.com/2011/07/05/forwarding-snmp-ports-over-ssh-using-socat/)
- https://blog.openziti.io/introducing-zrok
- https://blog.rom1v.com/2017/03/introducing-gnirehtet/

- [ ] Consider using IP_FREEBIND sockopt
      - https://www.freedesktop.org/wiki/Software/systemd/NetworkTarget

- [ ] https://tls-anvil.com/docs/Quick-Start/index

- [ ] Commonalize connector behaviours
- [ ] Use cancelable contexts throughout
- [ ] gRPC
- [ ] WSS
- [ ] XMPP
- [ ] ZMQ
- [ ] [nostr](https://github.com/nostr-protocol/nostr)

- (?) eBPF
- (?) Encode packet with protocol buffers/MessagePack/Apache Avro
- (?) Wrap [libevent](https://libevent.org) or use syscalls
- (?) Routing matrix
- (?) Replace handler functions with channels
- (?) Remove dependency on uhppoted-lib and uhppote-core
- [ ] httpd connector logo
      - https://graphicdesign.stackexchange.com/questions/159149/how-to-draw-parallel-inclined-surfaces-in-perspective

