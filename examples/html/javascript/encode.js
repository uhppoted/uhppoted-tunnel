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
