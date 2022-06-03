package tunnel

import (
	"reflect"
	"testing"
)

func TestPacketize(t *testing.T) {
	id := uint32(12345)
	msg := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}
	expected := []byte{0x00, 0x08, 0x00, 0x00, 0x30, 0x39, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}

	packet := packetize(id, msg)

	if !reflect.DeepEqual(packet, expected) {
		t.Errorf("Incorrect packet\n   expected:%#v\n   got:     %#v", expected, packet)
	}
}

func TestDepacketize(t *testing.T) {
	buffer := []byte{0x00, 0x08, 0x00, 0x00, 0x30, 0x39, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 'A', 'B', 'C', 'D'}
	expected := message{
		id:      12345,
		message: []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef},
	}

	id, msg, remaining := depacketize(buffer)

	if id != expected.id {
		t.Errorf("depacketize - incorrect ID, expected:%v, got:%v", expected.id, id)
	}

	if !reflect.DeepEqual(msg, expected.message) {
		t.Errorf("Incorrect packet\n   expected:%#v\n   got:     %#v", expected.message, msg)
	}

	if !reflect.DeepEqual(remaining, []byte{'A', 'B', 'C', 'D'}) {
		t.Errorf("Incorrect remaining\n   expected:%#v\n   got:     %#v", []byte{'A', 'B', 'C', 'D'}, remaining)
	}
}

func TestDepacketizeWithMissingHeader(t *testing.T) {
	buffer := []byte{0x00, 0x08, 0x00, 0x00, 0x30}

	id, msg, remaining := depacketize(buffer)

	if id != 0 {
		t.Errorf("depacketize - incorrect ID, expected:%v, got:%v", 0, id)
	}

	if msg != nil {
		t.Errorf("Incorrect packet\n   expected:%#v\n   got:     %#v", nil, msg)
	}

	if !reflect.DeepEqual(remaining, []byte{}) {
		t.Errorf("Incorrect remaining\n   expected:%#v\n   got:     %#v", []byte{}, remaining)
	}
}

func TestDepacketizeWithPartialMessage(t *testing.T) {
	buffer := []byte{0x00, 0x08, 0x00, 0x00, 0x30, 0x39, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd}

	id, msg, remaining := depacketize(buffer)

	if id != 0 {
		t.Errorf("depacketize - incorrect ID, expected:%v, got:%v", 0, id)
	}

	if msg != nil {
		t.Errorf("Incorrect packet\n   expected:%#v\n   got:     %#v", nil, msg)
	}

	if !reflect.DeepEqual(remaining, []byte{}) {
		t.Errorf("Incorrect remaining\n   expected:%#v\n   got:     %#v", []byte{}, remaining)
	}
}

func TestDepacketizeWithMultipleMessages(t *testing.T) {
	buffer := []byte{
		0x00, 0x08, 0x00, 0x00, 0x30, 0x39, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
		0x00, 0x08, 0x00, 0x00, 0x30, 0x3a, 0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10,
	}

	expected := []message{
		{id: 12345, message: []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}},
		{id: 12346, message: []byte{0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10}},
	}

	id, msg, buffer := depacketize(buffer)

	if id != expected[0].id {
		t.Errorf("depacketize - incorrect ID, expected:%v, got:%v", expected[0].id, id)
	}

	if !reflect.DeepEqual(msg, expected[0].message) {
		t.Errorf("Incorrect packet\n   expected:%#v\n   got:     %#v", expected[0].message, msg)
	}

	id, msg, buffer = depacketize(buffer)

	if id != expected[1].id {
		t.Errorf("depacketize - incorrect ID, expected:%v, got:%v", expected[1].id, id)
	}

	if !reflect.DeepEqual(msg, expected[1].message) {
		t.Errorf("Incorrect packet\n   expected:%#v\n   got:     %#v", expected[1].message, msg)
	}
}
