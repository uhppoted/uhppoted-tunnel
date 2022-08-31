import * as uhppote from './uhppote.js'

export function exec (fn, ...args) {
  return fn(...args)
}

export function getAllControllers () {
  return uhppote.GetAllControllers()
}

export function getController (controller) {
  return uhppote.GetController(
    document.querySelector(`input#${controller}`).value
  )
}

export function setIP (controller, address, netmask, gateway) {
  return uhppote.SetIP(
    document.querySelector(`input#${controller}`).value,
    document.querySelector(`input#${address}`).value,
    document.querySelector(`input#${netmask}`).value,
    document.querySelector(`input#${gateway}`).value
  )
}

export function getTime (controller) {
  return uhppote.GetTime(
    document.querySelector(`input#${controller}`).value
  )
}

export function setTime (controller, datetime) {
  const deviceID = document.querySelector(`input#${controller}`).value
  const dt = document.querySelector(`input#${datetime}`).value

  if (dt === '') {
    return uhppote.SetTime(deviceID, new Date())
  } else {
    return uhppote.SetTime(deviceID, new Date(datetime))
  }
}
