package tunnel

import (
	"sync/atomic"
)

var PACKETID uint32 = 0

func nextID() uint32 {
	return atomic.AddUint32(&PACKETID, 1)
}

func packetize(id uint32, message []byte) []byte {
	packet := make([]byte, len(message)+6)

	packet[0] = byte((len(message) >> 8) & 0x00ff)
	packet[1] = byte((len(message) >> 0) & 0x00ff)

	packet[2] = byte((id >> 24) & 0x00ff)
	packet[3] = byte((id >> 16) & 0x00ff)
	packet[4] = byte((id >> 8) & 0x00ff)
	packet[5] = byte((id >> 0) & 0x00ff)

	copy(packet[6:], message)

	return packet
}

func depacketize(packet []byte) (uint32, []byte) {
	if len(packet) < 6 {
		warnf("invalid packet (%v bytes)", len(packet))
	} else {
		N := int(packet[0])
		N <<= 8
		N += int(packet[1])

		id := uint32(packet[2])
		id <<= 8
		id += uint32(packet[3])
		id <<= 8
		id += uint32(packet[4])
		id <<= 8
		id += uint32(packet[5])

		if N > len(packet[6:]) {
			warnf("invalid packet - expected %v bytes, got %v bytes", N+6, len(packet))
		} else {
			return id, packet[6 : 6+N]
		}
	}

	return 0, nil
}
