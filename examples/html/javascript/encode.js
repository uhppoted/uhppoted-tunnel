export function GetDevices () {
  const request = new Uint8Array(64)

  request[0] = 0x17
  request[1] = 0x94

  return request
}

export function GetDevice (deviceID) {
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
