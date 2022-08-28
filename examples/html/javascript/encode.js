export function GetController (deviceID) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x94

  view.setUint32(4, deviceID, true)

  return request
}

export function SetIP (deviceID, address, netmask, gateway) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x96

  view.setUint32(4, deviceID, true)
  request.set(IPv4(address), 8)
  request.set(IPv4(netmask), 12)
  request.set(IPv4(gateway), 16)
  view.setUint32(20, 0x55aaaa55, true)

  return request
}

export function GetTime (deviceID) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)

  request[0] = 0x17
  request[1] = 0x32

  view.setUint32(4, deviceID, true)

  return request
}

export function SetTime (deviceID, datetime) {
  const request = new Uint8Array(64)
  const view = new DataView(request.buffer)
  const now = new Date()

  request[0] = 0x17
  request[1] = 0x30

  view.setUint32(4, deviceID, true)

  if (datetime === '') {
    request.set(datetime2bin(now), 8)
  } else {
    request.set(datetime2bin(new Date(datetime)), 8)
  }

  return request
}

function IPv4 (s) {
  const re = /([0-9]{0,3})\.([0-9]{0,3})\.([0-9]{0,3})\.([0-9]{0,3})/
  const match = s.match(re)
  const ip = []

  if (!match || match.length !== 5) {
    throw new Error(`invalid IP address ${s}`)
  }

  for (let i = 0; i < 4; i++) {
    const b = Number(match[i + 1])
    if (Number.isNaN(b) || b > 255) {
      throw new Error(`invalid IP address ${s}`)
    } else {
      ip.push(b)
    }
  }

  return ip
}

function datetime2bin (datetime) {
  const year = String(datetime.getFullYear()).padStart(4, '0')
  const month = String(datetime.getMonth() + 1).padStart(2, '0')
  const day = String(datetime.getDate()).padStart(2, '0')
  const hour = String(datetime.getHours()).padStart(2, '0')
  const minute = String(datetime.getMinutes()).padStart(2, '0')
  const second = String(datetime.getSeconds()).padStart(2, '0')

  const date = `${year}${month}${day}`
  const time = `${hour}${minute}${second}`

  return bcd2bin(`${date}${time}`)
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
