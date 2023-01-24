# TODO

## IN PROGRESS

- [x] TLS client to non-TLS host handshake doesn't ever timeout/disconnect
      - (?) Maybe do the same handshake as the server
```
./bin/uhppoted-tunnel --debug --console --in udp/event:0.0.0.0:60001 --out tls/client:127.0.0.1:12345
./bin/uhppoted-tunnel --debug --console --in tcp/server:0.0.0.0:12345 --out udp/event:192.168.1.255:60005
```

- [ ] Error on close
```
2023/01/18 11:22:55 INFO   ROUTER closing
2023/01/18 11:22:55 FATAL         runtime error: invalid memory address or nil pointer dereference
```



## TODO

- [ ] integration tests
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

