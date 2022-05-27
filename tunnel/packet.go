package tunnel

func packetize(message []byte) []byte {
	packet := make([]byte, len(message)+2)

	packet[0] = byte((len(message) >> 8) & 0x00ff)
	packet[1] = byte((len(message) >> 0) & 0x00ff)
	copy(packet[2:], message)

	return packet
}

func depacketize(packet []byte) []byte {
	return packet[2:]
}
