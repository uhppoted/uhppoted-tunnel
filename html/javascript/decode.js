export function GetDevice (bytes) {
  const buffer = new Uint8Array(bytes)
  const view = new DataView(buffer.buffer)

  return {
    device: {
      id: view.getUint32(4, true),
      address: address(view.getUint32(8)),
      netmask: address(view.getUint32(12)),
      gateway: address(view.getUint32(16)),
      MAC: MAC(buffer.slice(20, 26)),
      version: bcd(buffer.slice(26, 28)),
      date: yyyymmdd(buffer.slice(28, 32))
    }
  }
}

function address (ulong) {
  const b1 = (ulong >>> 24 & 0xff)
  const b2 = (ulong >>> 16 & 0xff)
  const b3 = (ulong >>> 8 & 0xff)
  const b4 = (ulong >>> 0 & 0xff)

  return `${b1}.${b2}.${b3}.${b4}`
}

function MAC (bytes) {
  return [...bytes].map(x => x.toString(16).padStart(2, '0')).join(' ')
}

function bcd (bytes) {
  return [...bytes].map(x => [(x >>> 4) & 0x0f, (x >>> 0) & 0x0f]).flat().join('')
}

function yyyymmdd (bytes) {
  const date = bcd(bytes)

  if (date === '00000000') {
    return ''
  }

  return date.substr(0, 4) + '-' + date.substr(4, 2) + '-' + date.substr(6, 2)
}
