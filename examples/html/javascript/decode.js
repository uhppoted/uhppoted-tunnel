export function decode (bytes) {
  const buffer = new Uint8Array(bytes)

  if (buffer.length !== 64) {
    throw new Error(`Invalid buffer (expected 64 bytes, got ${buffer.length}`)
  }

  if (buffer[0] !== 0x17) {
    throw new Error(`Invalid SOM ${buffer[0]}`)
  }

  switch (buffer[1]) {
    case 0x94:
      return GetDevice(bytes)

    case 0x32:
      return GetTime(bytes)

    default:
      throw new Error(`Unknown function code ${buffer[1]}`)
  }
}

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

export function GetTime (bytes) {
  const buffer = new Uint8Array(bytes)
  const view = new DataView(buffer.buffer)

  return {
    time: {
      id: view.getUint32(4, true),
      datetime: yyyymmddHHmmss(buffer.slice(8, 15))
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

function yyyymmddHHmmss (bytes) {
  const datetime = bcd(bytes)

  if (datetime === '00000000000000') {
    return ''
  }

  const date = `${datetime.substr(0, 4)}-${datetime.substr(4, 2)}-${datetime.substr(6, 2)}`
  const time = `${datetime.substr(8, 2)}:${datetime.substr(10, 2)}:${datetime.substr(12, 2)}`

  return `${date} ${time}`
}
