# TODO

## IN PROGRESS

- (??) DERP
       - https://pkg.go.dev/tailscale.com@v1.38.4/derp

- [ ] Rate limit (for e.g. when you've misconfigured your UDP broadcast to send to your own listener)

- [ ] Tailscale
      - [x] addr
      - [x] workdir
      - [x] hostname
      - [x] log format
      - [x] tsnet logging 
      - [ ] authorisation
            - [x] TS_AUTHKEY
            - [x] manual
            - [ ] TOML
                  - [x] authkey:
                  - [ ] ENV:
                  - [ ] OAuth2

      - [ ] plugin
      - [ ] README
            - [ ] TS_AUTHKEY
            - [ ] ...;debug for tailscale logging
            - [ ] NOTE: Keep alive is expected to be an application thing because if the server is 
                  down for a long time an internal keep-alive could backoff to the maximum limit
                  and then the next request could take a long time. 
      - [ ] CHANGELOG

      - https://tailscale.com/kb/1244/tsnet/
      - https://tailscale.dev/blog/embedded-funnel
      - https://github.com/tailscale/tailscale/issues/1552

- [ ] Error on CTRL-C
```
2023/04/19 10:07:03 FATAL         runtime error: invalid memory address or nil pointer dereference
```

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

