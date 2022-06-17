package protocol

import (
	"reflect"
	"testing"
)

func TestPacketize(t *testing.T) {
	id := uint32(12345)
	msg := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}
	expected := []byte{0x00, 0x08, 0x00, 0x00, 0x30, 0x39, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}

	packet := Packetize(id, msg)

	if !reflect.DeepEqual(packet, expected) {
		t.Errorf("Incorrect packet\n   expected:%#v\n   got:     %#v", expected, packet)
	}
}

func TestDepacketize(t *testing.T) {
	buffer := []byte{0x00, 0x08, 0x00, 0x00, 0x30, 0x39, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 'A', 'B', 'C', 'D'}
	expected := Message{
		ID:      12345,
		Message: []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef},
	}

	id, msg, remaining := Depacketize(buffer)

	if id != expected.ID {
		t.Errorf("depacketize - incorrect ID, expected:%v, got:%v", expected.ID, id)
	}

	if !reflect.DeepEqual(msg, expected.Message) {
		t.Errorf("Incorrect packet\n   expected:%#v\n   got:     %#v", expected.Message, msg)
	}

	if !reflect.DeepEqual(remaining, []byte{'A', 'B', 'C', 'D'}) {
		t.Errorf("Incorrect remaining\n   expected:%#v\n   got:     %#v", []byte{'A', 'B', 'C', 'D'}, remaining)
	}
}

func TestDepacketizeWithMissingHeader(t *testing.T) {
	buffer := []byte{0x00, 0x08, 0x00, 0x00, 0x30}

	id, msg, remaining := Depacketize(buffer)

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

	id, msg, remaining := Depacketize(buffer)

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

	expected := []Message{
		{ID: 12345, Message: []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}},
		{ID: 12346, Message: []byte{0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10}},
	}

	id, msg, buffer := Depacketize(buffer)

	if id != expected[0].ID {
		t.Errorf("depacketize - incorrect ID, expected:%v, got:%v", expected[0].ID, id)
	}

	if !reflect.DeepEqual(msg, expected[0].Message) {
		t.Errorf("Incorrect packet\n   expected:%#v\n   got:     %#v", expected[0].Message, msg)
	}

	id, msg, buffer = Depacketize(buffer)

	if id != expected[1].ID {
		t.Errorf("depacketize - incorrect ID, expected:%v, got:%v", expected[1].ID, id)
	}

	if !reflect.DeepEqual(msg, expected[1].Message) {
		t.Errorf("Incorrect packet\n   expected:%#v\n   got:     %#v", expected[1].Message, msg)
	}
}
