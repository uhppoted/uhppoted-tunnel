# TODO

- [x] Reverse tunnel
- [x] `depacketize` should return the remaining buffer
- [x] Cleanup TCPClient on broken pipe
- [x] Check works with multiple replies
- [ ] reverse-host/client: not always getting all replies

- [ ] Move handlers back to TCP/UDP and just use s.relay in switch
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
