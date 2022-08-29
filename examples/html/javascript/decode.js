
export function GetControllerResponse (packet) {
  const buffer = new Uint8Array(packet)
  const view = new DataView(buffer.buffer)

  if (buffer.length !== 64) {
    throw new Error(`invalid reply packet length (${buffer.length})`)
  }

  if (buffer[1] !== 0x94) {
    throw new Error(`invalid reply function code (${buffer[1].toString(16).padStart(2, '0')})`)
  }

  return {
    deviceId: unpackUint32(view, 4),
    ipAddress: unpackIPv4(view, 8),
    subnetMask: unpackIPv4(view, 12),
    gateway: unpackIPv4(view, 16),
    MACAddress: unpackMAC(view, 20),
    version: unpackVersion(view, 26),
    date: unpackDate(view, 28)
  }
}

export function GetTimeResponse (packet) {
  const buffer = new Uint8Array(packet)
  const view = new DataView(buffer.buffer)

  if (buffer.length !== 64) {
    throw new Error(`invalid reply packet length (${buffer.length})`)
  }

  if (buffer[1] !== 0x32) {
    throw new Error(`invalid reply function code (${buffer[1].toString(16).padStart(2, '0')})`)
  }

  return {
    deviceId: unpackUint32(view, 4),
    datetime: unpackDatetime(view, 8)
  }
}

export function SetTimeResponse (packet) {
  const buffer = new Uint8Array(packet)
  const view = new DataView(buffer.buffer)

  if (buffer.length !== 64) {
    throw new Error(`invalid reply packet length (${buffer.length})`)
  }

  if (buffer[1] !== 0x30) {
    throw new Error(`invalid reply function code (${buffer[1].toString(16).padStart(2, '0')})`)
  }

  return {
    deviceId: unpackUint32(view, 4),
    datetime: unpackDatetime(view, 8)
  }
}

export function GetStatusResponse (packet) {
  const buffer = new Uint8Array(packet)
  const view = new DataView(buffer.buffer)

  if (buffer.length !== 64) {
    throw new Error(`invalid reply packet length (${buffer.length})`)
  }

  if (buffer[1] !== 0x20) {
    throw new Error(`invalid reply function code (${buffer[1].toString(16).padStart(2, '0')})`)
  }

  return {
    deviceId: unpackUint32(view, 4),
    systemDate: unpackShortdate(view, 51),
    systemTime: unpackTime(view, 37),
    door1Open: unpackBool(view, 28),
    door2Open: unpackBool(view, 29),
    door3Open: unpackBool(view, 30),
    door4Open: unpackBool(view, 31),
    door1Button: unpackBool(view, 32),
    door2Button: unpackBool(view, 33),
    door3Button: unpackBool(view, 34),
    door4Button: unpackBool(view, 35),
    relays: unpackUint8(view, 49),
    inputs: unpackUint8(view, 50),
    systemError: unpackUint8(view, 36),
    specialInfo: unpackUint8(view, 48),
    eventIndex: unpackUint32(view, 8),
    eventType: unpackUint8(view, 12),
    eventAccessGranted: unpackBool(view, 13),
    eventDoor: unpackUint8(view, 14),
    eventDirection: unpackUint8(view, 15),
    eventCard: unpackUint32(view, 16),
    eventTimestamp: unpackDatetime(view, 20),
    eventReason: unpackUint8(view, 27),
    sequenceNo: unpackUint32(view, 40)
  }
}

export function GetListenerResponse (packet) {
  const buffer = new Uint8Array(packet)
  const view = new DataView(buffer.buffer)

  if (buffer.length !== 64) {
    throw new Error(`invalid reply packet length (${buffer.length})`)
  }

  if (buffer[1] !== 0x92) {
    throw new Error(`invalid reply function code (${buffer[1].toString(16).padStart(2, '0')})`)
  }

  return {
    deviceId: unpackUint32(view, 4),
    address: unpackIPv4(view, 8),
    port: unpackUint16(view, 12)
  }
}

export function SetListenerResponse (packet) {
  const buffer = new Uint8Array(packet)
  const view = new DataView(buffer.buffer)

  if (buffer.length !== 64) {
    throw new Error(`invalid reply packet length (${buffer.length})`)
  }

  if (buffer[1] !== 0x90) {
    throw new Error(`invalid reply function code (${buffer[1].toString(16).padStart(2, '0')})`)
  }

  return {
    deviceId: unpackUint32(view, 4),
    ok: unpackBool(view, 8)
  }
}

export function GetDoorControlResponse (packet) {
  const buffer = new Uint8Array(packet)
  const view = new DataView(buffer.buffer)

  if (buffer.length !== 64) {
    throw new Error(`invalid reply packet length (${buffer.length})`)
  }

  if (buffer[1] !== 0x82) {
    throw new Error(`invalid reply function code (${buffer[1].toString(16).padStart(2, '0')})`)
  }

  return {
    deviceId: unpackUint32(view, 4),
    door: unpackUint8(view, 8),
    mode: unpackUint8(view, 9),
    delay: unpackUint8(view, 10)
  }
}

export function SetDoorControlResponse (packet) {
  const buffer = new Uint8Array(packet)
  const view = new DataView(buffer.buffer)

  if (buffer.length !== 64) {
    throw new Error(`invalid reply packet length (${buffer.length})`)
  }

  if (buffer[1] !== 0x80) {
    throw new Error(`invalid reply function code (${buffer[1].toString(16).padStart(2, '0')})`)
  }

  return {
    deviceId: unpackUint32(view, 4),
    door: unpackUint8(view, 8),
    mode: unpackUint8(view, 9),
    delay: unpackUint8(view, 10)
  }
}

export function OpenDoorResponse (packet) {
  const buffer = new Uint8Array(packet)
  const view = new DataView(buffer.buffer)

  if (buffer.length !== 64) {
    throw new Error(`invalid reply packet length (${buffer.length})`)
  }

  if (buffer[1] !== 0x40) {
    throw new Error(`invalid reply function code (${buffer[1].toString(16).padStart(2, '0')})`)
  }

  return {
    deviceId: unpackUint32(view, 4),
    opened: unpackBool(view, 8)
  }
}

export function GetCardsResponse (packet) {
  const buffer = new Uint8Array(packet)
  const view = new DataView(buffer.buffer)

  if (buffer.length !== 64) {
    throw new Error(`invalid reply packet length (${buffer.length})`)
  }

  if (buffer[1] !== 0x58) {
    throw new Error(`invalid reply function code (${buffer[1].toString(16).padStart(2, '0')})`)
  }

  return {
    deviceId: unpackUint32(view, 4),
    cards: unpackUint32(view, 8)
  }
}

export function GetCardResponse (packet) {
  const buffer = new Uint8Array(packet)
  const view = new DataView(buffer.buffer)

  if (buffer.length !== 64) {
    throw new Error(`invalid reply packet length (${buffer.length})`)
  }

  if (buffer[1] !== 0x5a) {
    throw new Error(`invalid reply function code (${buffer[1].toString(16).padStart(2, '0')})`)
  }

  return {
    deviceId: unpackUint32(view, 4),
    cardNumber: unpackUint32(view, 8),
    startDate: unpackDate(view, 12),
    endDate: unpackDate(view, 16),
    door1: unpackUint8(view, 20),
    door2: unpackUint8(view, 21),
    door3: unpackUint8(view, 22),
    door4: unpackUint8(view, 23)
  }
}

export function GetCardByIndexResponse (packet) {
  const buffer = new Uint8Array(packet)
  const view = new DataView(buffer.buffer)

  if (buffer.length !== 64) {
    throw new Error(`invalid reply packet length (${buffer.length})`)
  }

  if (buffer[1] !== 0x5c) {
    throw new Error(`invalid reply function code (${buffer[1].toString(16).padStart(2, '0')})`)
  }

  return {
    deviceId: unpackUint32(view, 4),
    cardNumber: unpackUint32(view, 8),
    startDate: unpackDate(view, 12),
    endDate: unpackDate(view, 16),
    door1: unpackUint8(view, 20),
    door2: unpackUint8(view, 21),
    door3: unpackUint8(view, 22),
    door4: unpackUint8(view, 23)
  }
}

export function PutCardResponse (packet) {
  const buffer = new Uint8Array(packet)
  const view = new DataView(buffer.buffer)

  if (buffer.length !== 64) {
    throw new Error(`invalid reply packet length (${buffer.length})`)
  }

  if (buffer[1] !== 0x50) {
    throw new Error(`invalid reply function code (${buffer[1].toString(16).padStart(2, '0')})`)
  }

  return {
    deviceId: unpackUint32(view, 4),
    stored: unpackBool(view, 8)
  }
}

export function DeleteCardResponse (packet) {
  const buffer = new Uint8Array(packet)
  const view = new DataView(buffer.buffer)

  if (buffer.length !== 64) {
    throw new Error(`invalid reply packet length (${buffer.length})`)
  }

  if (buffer[1] !== 0x52) {
    throw new Error(`invalid reply function code (${buffer[1].toString(16).padStart(2, '0')})`)
  }

  return {
    deviceId: unpackUint32(view, 4),
    deleted: unpackBool(view, 8)
  }
}

export function DeleteAllCardsResponse (packet) {
  const buffer = new Uint8Array(packet)
  const view = new DataView(buffer.buffer)

  if (buffer.length !== 64) {
    throw new Error(`invalid reply packet length (${buffer.length})`)
  }

  if (buffer[1] !== 0x54) {
    throw new Error(`invalid reply function code (${buffer[1].toString(16).padStart(2, '0')})`)
  }

  return {
    deviceId: unpackUint32(view, 4),
    deleted: unpackBool(view, 8)
  }
}

export function GetEventResponse (packet) {
  const buffer = new Uint8Array(packet)
  const view = new DataView(buffer.buffer)

  if (buffer.length !== 64) {
    throw new Error(`invalid reply packet length (${buffer.length})`)
  }

  if (buffer[1] !== 0xb0) {
    throw new Error(`invalid reply function code (${buffer[1].toString(16).padStart(2, '0')})`)
  }

  return {
    deviceId: unpackUint32(view, 4),
    index: unpackUint32(view, 8),
    eventType: unpackUint8(view, 12),
    accessGranted: unpackBool(view, 13),
    door: unpackUint8(view, 14),
    direction: unpackUint8(view, 15),
    card: unpackUint32(view, 16),
    timestamp: unpackDatetime(view, 20),
    reason: unpackUint8(view, 27)
  }
}

export function GetEventIndexResponse (packet) {
  const buffer = new Uint8Array(packet)
  const view = new DataView(buffer.buffer)

  if (buffer.length !== 64) {
    throw new Error(`invalid reply packet length (${buffer.length})`)
  }

  if (buffer[1] !== 0xb4) {
    throw new Error(`invalid reply function code (${buffer[1].toString(16).padStart(2, '0')})`)
  }

  return {
    deviceId: unpackUint32(view, 4),
    eventIndex: unpackUint32(view, 8)
  }
}

export function SetEventIndexResponse (packet) {
  const buffer = new Uint8Array(packet)
  const view = new DataView(buffer.buffer)

  if (buffer.length !== 64) {
    throw new Error(`invalid reply packet length (${buffer.length})`)
  }

  if (buffer[1] !== 0xb2) {
    throw new Error(`invalid reply function code (${buffer[1].toString(16).padStart(2, '0')})`)
  }

  return {
    deviceId: unpackUint32(view, 4),
    updated: unpackBool(view, 8)
  }
}

export function RecordSpecialEventsResponse (packet) {
  const buffer = new Uint8Array(packet)
  const view = new DataView(buffer.buffer)

  if (buffer.length !== 64) {
    throw new Error(`invalid reply packet length (${buffer.length})`)
  }

  if (buffer[1] !== 0x8e) {
    throw new Error(`invalid reply function code (${buffer[1].toString(16).padStart(2, '0')})`)
  }

  return {
    deviceId: unpackUint32(view, 4),
    enabled: unpackBool(view, 8)
  }
}

export function GetTimeProfileResponse (packet) {
  const buffer = new Uint8Array(packet)
  const view = new DataView(buffer.buffer)

  if (buffer.length !== 64) {
    throw new Error(`invalid reply packet length (${buffer.length})`)
  }

  if (buffer[1] !== 0x98) {
    throw new Error(`invalid reply function code (${buffer[1].toString(16).padStart(2, '0')})`)
  }

  return {
    deviceId: unpackUint32(view, 4),
    profileId: unpackUint8(view, 8),
    startDate: unpackDate(view, 9),
    endDate: unpackDate(view, 13),
    monday: unpackBool(view, 17),
    tuesday: unpackBool(view, 18),
    wednesday: unpackBool(view, 19),
    thursday: unpackBool(view, 20),
    friday: unpackBool(view, 21),
    saturday: unpackBool(view, 22),
    sunday: unpackBool(view, 23),
    segment1Start: unpackHHmm(view, 24),
    segment1End: unpackHHmm(view, 26),
    segment2Start: unpackHHmm(view, 28),
    segment2End: unpackHHmm(view, 30),
    segment3Start: unpackHHmm(view, 32),
    segment3End: unpackHHmm(view, 34),
    linkedProfileId: unpackUint8(view, 36)
  }
}

export function SetTimeProfileResponse (packet) {
  const buffer = new Uint8Array(packet)
  const view = new DataView(buffer.buffer)

  if (buffer.length !== 64) {
    throw new Error(`invalid reply packet length (${buffer.length})`)
  }

  if (buffer[1] !== 0x88) {
    throw new Error(`invalid reply function code (${buffer[1].toString(16).padStart(2, '0')})`)
  }

  return {
    deviceId: unpackUint32(view, 4),
    stored: unpackBool(view, 8)
  }
}

export function DeleteAllTimeProfilesResponse (packet) {
  const buffer = new Uint8Array(packet)
  const view = new DataView(buffer.buffer)

  if (buffer.length !== 64) {
    throw new Error(`invalid reply packet length (${buffer.length})`)
  }

  if (buffer[1] !== 0x8a) {
    throw new Error(`invalid reply function code (${buffer[1].toString(16).padStart(2, '0')})`)
  }

  return {
    deviceId: unpackUint32(view, 4),
    deleted: unpackBool(view, 8)
  }
}

function unpackUint8 (packet, offset) {
  return packet.getUint8(offset)
}

function unpackUint16 (packet, offset) {
  return packet.getUint16(offset, true)
}

function unpackUint32 (packet, offset) {
  return packet.getUint32(offset, true)
}

function unpackBool (packet, offset) {
  return packet.getUint8(offset) !== 0x00
}

function unpackIPv4 (packet, offset) {
  const bytes = new Uint8Array(packet.buffer.slice(offset, offset + 4))

  return [...bytes].map(x => x.toString(10)).join('.')
}

function unpackMAC (packet, offset) {
  const bytes = new Uint8Array(packet.buffer.slice(offset, offset + 6))

  return [...bytes].map(x => x.toString(16).padStart(2, '0')).join(':')
}

function unpackVersion (packet, offset) {
  const major = packet.getUint8(offset).toString(16)
  const minor = packet.getUint8(offset + 1).toString(16).padStart(2, '0')

  return `v${major}.${minor}`
}

function unpackDate (packet, offset) {
  const bytes = new Uint8Array(packet.buffer.slice(offset, offset + 4))
  const datetime = bcd(bytes)

  if (datetime === '00000000') {
    return ''
  }

  const date = `${datetime.substr(0, 4)}-${datetime.substr(4, 2)}-${datetime.substr(6, 2)}`

  return `${date}`
}

function unpackShortdate (packet, offset) {
  const bytes = new Uint8Array(packet.buffer.slice(offset, offset + 3))
  const datetime = bcd('20' + bytes)

  if (datetime === '20000000') {
    return ''
  }

  const date = `${datetime.substr(0, 4)}-${datetime.substr(4, 2)}-${datetime.substr(6, 2)}`

  return `${date}`
}

function unpackDatetime (packet, offset) {
  const bytes = new Uint8Array(packet.buffer.slice(offset, offset + 7))
  const datetime = bcd(bytes)

  if (datetime === '00000000000000') {
    return ''
  }

  const date = `${datetime.substr(0, 4)}-${datetime.substr(4, 2)}-${datetime.substr(6, 2)}`
  const time = `${datetime.substr(8, 2)}:${datetime.substr(10, 2)}:${datetime.substr(12, 2)}`

  return `${date} ${time}`
}

function unpackTime (packet, offset) {
  const bytes = new Uint8Array(packet.buffer.slice(offset, offset + 3))
  const datetime = bcd(bytes)

  if (datetime === '000000') {
    return ''
  }

  const time = `${datetime.substr(8, 2)}:${datetime.substr(10, 2)}:${datetime.substr(12, 2)}`

  return `${time}`
}

function unpackHHmm (packet, offset) {
  const bytes = new Uint8Array(packet.buffer.slice(offset, offset + 2))
  const datetime = bcd(bytes)

  if (datetime === '0000') {
    return ''
  }

  const time = `${datetime.substr(8, 2)}:${datetime.substr(10, 2)}}`

  return `${time}`
}

function bcd (bytes) {
  return [...bytes].map(x => [(x >>> 4) & 0x0f, (x >>> 0) & 0x0f]).flat().join('')
}
