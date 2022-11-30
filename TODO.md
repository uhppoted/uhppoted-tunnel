# TODO

## IN PROGRESS

- [ ] https://github.com/uhppoted/uhppoted-tunnel/issues/7
      - [x] Move lockfile implementation to uhppoted-lib
      - [x] Default to ephemeral _tmp_ folder for lockfiles
      - [ ] Use `flock` for default implementation 
            - [x] Darwin
            - [x] Linux
            - [ ] Windows
            - https://stackoverflow.com/questions/34710460/golang-flock-filelocking-throwing-panic-runtime-error-invalid-memory-address-o
            - https://stackoverflow.com/questions/52986413/how-to-get-an-exclusive-lock-on-a-file-in-go
            - https://github.com/gofrs/flock
      - [ ] Clean up flocked lockfile
            - https://stackoverflow.com/questions/17708885/flock-removing-locked-file-without-race-condition
            - (?) unlink
      - [ ] Replace soft lock (MQTT) with flock
      - (?) Use cancelable context to release lock
      - [ ] Figure out what on earth this thinks it is doing ??????
```
      defer func() {
            if err := recover(); err != nil {
                  fatalf("%v", err)
            }
      }()
```
            - https://stackoverflow.com/questions/34710460/golang-flock-filelocking-throwing-panic-runtime-error-invalid-memory-address-o
            - http://blog.golang.org/defer-panic-and-recover

- [x] ARM64 build
      - [ ] Test on Google VM

- [ ] log.Warnf+ should default to stderr
- [ ] Windows eventlog message file
      - https://social.msdn.microsoft.com/Forums/windowsdesktop/en-US/deaa0055-7770-4e55-a5b8-6d08b80b74af/creating-event-log-message-files
      - FormatMessage
         - https://go.dev/src/syscall/syscall_windows.go

## TODO

- (?) https://eli.thegreenplace.net/2022/ssh-port-forwarding-with-go/
- (?) [UDP tunnelling: ssh/nc](https://superuser.com/questions/53103/udp-traffic-through-ssh-tunnel)
- (?) [UDP tunnelling: socat](http://www.morch.com/2011/07/05/forwarding-snmp-ports-over-ssh-using-socat/)

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

