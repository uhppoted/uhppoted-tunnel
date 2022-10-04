# TODO

## IN PROGRESS

- [x] https://github.com/uhppoted/uhppoted-tunnel/issues/4
      - [x] Add mutex to Router::get
      - [x] Maybe abstract _handlers_ out into a mutexed map or something
      - [ ] Run extended test on VPS

- [x] Logo
      - [x] Add freepik attribution for HTTP/S logo.

## TODO

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

