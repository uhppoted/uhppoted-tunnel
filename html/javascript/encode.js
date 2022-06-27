export function GetDevices () {
  const request = new Uint8Array(64)

  request[0] = 0x17
  request[1] = 0x94

  return request
}
