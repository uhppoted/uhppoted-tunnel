# TODO

- [x] Reverse tunnel
- [x] `depacketize` should return the remaining buffer
- [x] Cleanup TCPClient on broken pipe
- [ ] Check works with multiple replies
```
 ... receive error: invalid device ID - expected:201020304, got:303986753
       201020304  192.168.1.101  255.255.255.0  192.168.1.1  52:fd:fc:07:21:82  v6.62  2020-01-01
Alpha  405419896  192.168.1.100  255.255.255.0  192.168.1.1  00:12:23:34:45:56  v8.92  2018-11-05
```

- [ ] Clean up switch handlers after timeout
- [ ] Routing matrix
- [ ] Replace handler functions with channels

## Miscellaneous
- [ ] Encode packet with protocol buffers
- [ ] lockfiles
      - command line option
- [ ] daemonize
- [ ] undaemonize
- [ ] SSL + mutual auth
- [ ] redirect events
- (?) Wrap [libevent](https://libevent.org)
