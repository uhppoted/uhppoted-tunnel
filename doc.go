// Copyright 2023 uhppoted@twyst.co.za. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

/*
Package uhppote-tunnel implements a relay/switch for the UDP packets that 
control the UHPPOTE TCP/IP Wiegand-26 access controllers.

The tunnel connectors are primarily designed to allow the access controllers 
to be managed remotely via a TCP pipe secured with mutually authenticated TLS,
but can be used for e.g:

- an HTTP CLI client
- event collation 
- event fanout

# Connectors

The package includes the following connectors:
  
  - udp/listen: receives and relays UDP commands from a management application and return the replies.
  - udp/broadcast: broadcasts UDP commands to the access controllers and relays the replies.
  - udp/event: relays access controller events
  - tcp/client: bidirectional TCP/IP pipe that connects to a remote server and relays commands and replies
  - tcp/server: bidirectional TCP/IP pipe that accepts remote connections and relays commands and replies
  - tls/client: tcp/client connector secured with TLS
  - tls/server: tcp/client connector secured with TLS
  - http: relays commands submitted as HTTP POST requests and returns the reply
*/
package tunnel
