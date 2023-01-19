# TODO

## IN PROGRESS

- [ ] 'events' connectors
      - [x] UDP event/in
      - [x] UDP event/out
      - [x] TCP event/out
      - [x] TCP event/in
      - [x] TLS event/out
      - [x] TLS event/in
      - [x] Reverse tunnel
      - [ ] Clean up in/out connector semantics
      - [ ] Restructure tunnel connectors so that you are composable/constructable
      - [ ] README

- [ ] TLS client to non-TLS host handshake doesn't ever timeout/disconnect
      - (?) Maybe do the same handshake as the server

- [ ] Error on close
```
2023/01/18 11:22:55 INFO   ROUTER closing
2023/01/18 11:22:55 FATAL         runtime error: invalid memory address or nil pointer dereference
```

- [ ] Maybe don't send 0 length packets ?
      - Happens e.g. when a TLS client tries to connect to a TCP server
```
./bin/uhppoted-tunnel --debug --console --in tcp/event:0.0.0.0:12345 --out udp/event:192.168.1.255:60005 --udp-timeout 1s
...
...
2023/01/18 12:06:55 DEBUG  TCP    received 215 bytes from 127.0.0.1:50568
                                  00000000  16 03 01 00 d2 01 00 00  ce 03 03 0f 71 11 4f 96
                                  00000010  91 98 88 ea 4e 40 43 fe  32 c1 1f a2 a3 53 9b dd
                                  00000020  22 7a 5d 1e 89 68 08 e2  22 88 35 20 19 dc 07 92
                                  00000030  f6 e0 a2 a1 52 5b 23 f6  97 db b4 28 0b 96 0b 95
                                  00000040  ca bb ef 72 11 48 7c 64  0e dc 20 98 00 0e c0 2bd.. ....+|
                                  00000050  c0 2f c0 2c c0 30 13 01  13 02 13 03 01 00 00 77
                                  00000060  00 05 00 05 01 00 00 00  00 00 0a 00 0a 00 08 00
                                  00000070  1d 00 17 00 18 00 19 00  0b 00 02 01 00 00 0d 00
                                  00000080  1a 00 18 08 04 04 03 08  07 08 05 08 06 04 01 05
                                  00000090  01 06 01 05 03 06 03 02  01 02 03 ff 01 00 01 00
                                  000000a0  00 12 00 00 00 2b 00 05  04 03 04 03 03 00 33 00
                                  000000b0  26 00 24 00 1d 00 20 db  0f 95 dd e4 87 0a 5c fa
                                  000000c0  6f 1f 64 20 30 19 9b f7  88 40 f9 d9 47 14 c8 45
                                  000000d0  8c 04 7a 51 e2 60 26
                                  
2023/01/18 12:06:55 WARN          invalid packet - expected 5641 bytes, got 215 bytes
2023/01/18 12:06:55 DEBUG  UDP    event/out (0 bytes)
                                  
2023/01/18 12:06:55 DEBUG  UDP    sent 0 bytes to 192.168.1.255:60005
 ...
 ...
(uhppote-cli) ERROR: invalid message length - expected:64, got:0
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

