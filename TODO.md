# TODO

## IN PROGRESS

- [x] https://github.com/uhppoted/uhppoted-tunnel/issues/2
      - [x] Load TOML file
      - [x] Specify TOML file section on command line
      - [x] Use # or :: as seperator
      - [x] Override with command-line args
      - [x] Somehow roll lib.CommandX into lib.Command
      - [x] Clean up initial command configuration
      - [x] Merge back into _master_

- [ ] https://github.com/uhppoted/uhppoted-tunnel/issues/7
      - https://stackoverflow.com/questions/5210945/atomic-file-creation-on-linux
      - https://www.diskodev.com/posts/linux-atomic-operations-on-files/
      - https://stackoverflow.com/questions/29261648/atomic-writing-to-file-on-linux
      - https://stackoverflow.com/questions/34873151/how-can-i-delete-a-unix-domain-socket-file-when-i-exit-my-application
      - https://gavv.net/articles/unix-socket-reuse/
      - https://stackoverflow.com/questions/7405932/how-to-know-whether-any-process-is-bound-to-a-unix-domain-socket
      - https://www.unix.com/man-page/centos/1/lockfile/
      - https://security.stackexchange.com/questions/11976/can-unix-domain-sockets-be-locked-by-user-id
      - https://go.googlesource.com/proposal/+/master/design/33974-add-public-lockedfile-pkg.md
      - https://github.com/golang/go/issues/33974
      - https://pythonhosted.org/lockfile/lockfile.html

## TODO

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

