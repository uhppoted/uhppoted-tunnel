package tunnel

func packetize(message []byte) []byte {
	packet := make([]byte, len(message)+2)

	packet[0] = byte((len(message) >> 8) & 0x00ff)
	packet[1] = byte((len(message) >> 0) & 0x00ff)
	copy(packet[2:], message)

	return packet
}

func depacketize(packet []byte) []byte {
	if len(packet) < 2 {
		warnf("invalid packet (%v bytes)", len(packet))
	} else {
		N := int(packet[0])
		N <<= 8
		N += int(packet[1])

		if N > len(packet[2:]) {
			warnf("invalid packet - expected %v bytes, got %v bytes", N+2, len(packet))
		} else {
			return packet[2 : 2+N]
		}
	}

	return nil
}
