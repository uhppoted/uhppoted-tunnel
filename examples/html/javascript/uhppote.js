import * as encode from './encode.js'
import * as decode from './decode.js'
import * as udp from './udp.js'

export function GetAllControllers () {
  const bytes = encode.GetControllerRequest(0)

  return udp.post(bytes, '500ms')
    .then(replies => {
      const list = []

      for (const reply of replies) {
        list.push(decode.GetController(reply))
      }

      return list
    })
}

export function GetController (deviceID) {
  const bytes = encode.GetControllerRequest(deviceID)

  return udp.post(bytes, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.GetController(replies[0])
      }

      return null
    })
}

export function SetIP (deviceID, address, netmask, gateway) {
  const bytes = encode.SetIPRequest(deviceID, address, netmask, gateway)

  return udp.post(bytes, '0.1ms')
    .then(replies => {
      return true
    })
}

export function GetTime (deviceID) {
  const bytes = encode.GetTimeRequest(deviceID)

  return udp.post(bytes, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.GetTime(replies[0])
      }

      return null
    })
}

export function SetTime (deviceID, time) {
  const bytes = encode.SetTimeRequest(deviceID, time)

  return udp.post(bytes, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.GetTime(replies[0])
      }

      return null
    })
}
