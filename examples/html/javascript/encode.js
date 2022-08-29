
export function GetControllerRequest (deviceId) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x94

  packUint32(deviceId, view, 4)

  return request
}

export function SetIPRequest (deviceId, address, netmask, gateway) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x96

  packUint32(deviceId, view, 4)
  packIPv4(address, view, 8)
  packIPv4(netmask, view, 12)
  packIPv4(gateway, view, 16)
  packUint32(0x55aaaa55, view, 20)

  return request
}

export function GetTimeRequest (deviceId) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x32

  packUint32(deviceId, view, 4)

  return request
}

export function SetTimeRequest (deviceId, datetime) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x30

  packUint32(deviceId, view, 4)
  packDatetime(datetime, view, 8)

  return request
}

export function GetStatusRequest (deviceId) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x20

  packUint32(deviceId, view, 4)

  return request
}

export function GetListenerRequest (deviceId) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x92

  packUint32(deviceId, view, 4)

  return request
}

export function SetListenerRequest (deviceId, address, port) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x90

  packUint32(deviceId, view, 4)
  packIPv4(address, view, 8)
  packUint16(port, view, 12)

  return request
}

export function GetDoorControlRequest (deviceId, door) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x82

  packUint32(deviceId, view, 4)
  packUint8(door, view, 8)

  return request
}

export function SetDoorControlRequest (deviceId, door, mode, delay) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x80

  packUint32(deviceId, view, 4)
  packUint8(door, view, 8)
  packUint8(mode, view, 9)
  packUint8(delay, view, 10)

  return request
}

export function OpenDoorRequest (deviceId, door) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x40

  packUint32(deviceId, view, 4)
  packUint8(door, view, 8)

  return request
}

export function GetCardsRequest (deviceId) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x58

  packUint32(deviceId, view, 4)

  return request
}

export function GetCardRequest (deviceId, cardNumber) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x5a

  packUint32(deviceId, view, 4)
  packUint32(cardNumber, view, 8)

  return request
}

export function GetCardByIndexRequest (deviceId, cardIndex) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x5c

  packUint32(deviceId, view, 4)
  packUint32(cardIndex, view, 8)

  return request
}

export function PutCardRequest (deviceId, cardNumber, startDate, endDate, door1, door2, door3, door4) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x50

  packUint32(deviceId, view, 4)
  packUint32(cardNumber, view, 8)
  packDate(startDate, view, 12)
  packDate(endDate, view, 16)
  packUint8(door1, view, 20)
  packUint8(door2, view, 21)
  packUint8(door3, view, 22)
  packUint8(door4, view, 23)

  return request
}

export function DeleteCardRequest (deviceId, cardNumber) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x52

  packUint32(deviceId, view, 4)
  packUint32(cardNumber, view, 8)

  return request
}

export function DeleteCardsRequest (deviceId) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x54

  packUint32(deviceId, view, 4)
  packUint32(0x55aaaa55, view, 8)

  return request
}

export function GetEventRequest (deviceId, eventIndex) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0xb0

  packUint32(deviceId, view, 4)
  packUint32(eventIndex, view, 8)

  return request
}

export function GetEventIndexRequest (deviceId) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0xb4

  packUint32(deviceId, view, 4)

  return request
}

export function SetEventIndexRequest (deviceId, eventIndex) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0xb2

  packUint32(deviceId, view, 4)
  packUint32(eventIndex, view, 8)
  packUint32(0x55aaaa55, view, 12)

  return request
}

export function RecordSpecialEventsRequest (deviceId, enable) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x8e

  packUint32(deviceId, view, 4)
  packBool(enable, view, 8)

  return request
}

export function GetTimeProfileRequest (deviceId, profileId) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x98

  packUint32(deviceId, view, 4)
  packUint8(profileId, view, 8)

  return request
}

export function SetTimeProfileRequest (deviceId, profileId, startDate, endDate, monday, tuesday, wednesday, thursday, friday, saturday, sunday, segment1Start, segment1End, segment2Start, segment2End, segment3Start, segment3End, linkedProfileId) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x88

  packUint32(deviceId, view, 4)
  packUint8(profileId, view, 8)
  packDate(startDate, view, 9)
  packDate(endDate, view, 13)
  packBool(monday, view, 17)
  packBool(tuesday, view, 18)
  packBool(wednesday, view, 19)
  packBool(thursday, view, 20)
  packBool(friday, view, 21)
  packBool(saturday, view, 22)
  packBool(sunday, view, 23)
  packHHmm(segment1Start, view, 24)
  packHHmm(segment1End, view, 26)
  packHHmm(segment2Start, view, 28)
  packHHmm(segment2End, view, 30)
  packHHmm(segment3Start, view, 32)
  packHHmm(segment3End, view, 34)
  packUint8(linkedProfileId, view, 36)

  return request
}

export function DeleteAllTimeProfilesRequest (deviceId) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x8a

  packUint32(deviceId, view, 4)
  packUint32(0x55aaaa55, view, 8)

  return request
}

function packUint8 (v, packet, offset) {
  packet.setUint8(offset, v)
}

function packUint16 (v, packet, offset) {
  packet.setUint16(offset, v, true)
}

function packUint32 (v, packet, offset) {
  packet.setUint32(offset, v, true)
}

function packBool (v, packet, offset) {
  if (v) {
    packet.setUint8(offset, 0x01)
  } else {
    packet.setUint8(offset, 0x00)
  }
}

function packIPv4 (v, packet, offset) {
  const re = /([0-9]{0,3})\.([0-9]{0,3})\.([0-9]{0,3})\.([0-9]{0,3})/
  const match = v.match(re)

  if (!match || match.length !== 5) {
    throw new Error(`invalid IP address ${v}`)
  }

  for (let i = 0; i < 4; i++) {
    const b = Number(match[i + 1])
    if (Number.isNaN(b) || b > 255) {
      throw new Error(`invalid IP address ${v}`)
    } else {
      packet.setUint8(offset + i, b)
    }
  }
}

function packDate (v, packet, offset) {
  const year = String(v.getFullYear()).padStart(4, '0')
  const month = String(v.getMonth() + 1).padStart(2, '0')
  const day = String(v.getDate()).padStart(2, '0')

  const date = `${year}${month}${day}`
  const bytes = bcd2bin(`${date}`)

  for (let i = 0; i < 4; i++) {
    packet.setUint8(offset + i, bytes[i])
  }
}

function packDatetime (v, packet, offset) {
  const year = String(v.getFullYear()).padStart(4, '0')
  const month = String(v.getMonth() + 1).padStart(2, '0')
  const day = String(v.getDate()).padStart(2, '0')
  const hour = String(v.getHours()).padStart(2, '0')
  const minute = String(v.getMinutes()).padStart(2, '0')
  const second = String(v.getSeconds()).padStart(2, '0')

  const date = `${year}${month}${day}`
  const time = `${hour}${minute}${second}`
  const bytes = bcd2bin(`${date}${time}`)

  for (let i = 0; i < 7; i++) {
    packet.setUint8(offset + i, bytes[i])
  }
}

function packHHmm (v, packet, offset) {
  const hour = String(v.getHours()).padStart(2, '0')
  const minute = String(v.getMinutes()).padStart(2, '0')

  const time = `${hour}${minute}`
  const bytes = bcd2bin(`${time}`)

  packet.setUint8(offset, bytes[0])
  packet.setUint8(offset + 1, bytes[1])
}

function bcd2bin (bcd) {
  const bytes = []
  const matches = [...bcd.matchAll(/([0-9]{2})/g)]

  for (const m of matches) {
    const b = parseInt(m[0], 10)
    const byte = ((b / 10) << 4) | (b % 10)

    bytes.push(byte)
  }

  return bytes
}
