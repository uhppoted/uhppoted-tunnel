import * as uhppote from './uhppote.js'

export const commands = new Map([
  ['get-all-controllers', { fn: getAllControllers, args: [] }],
  ['get-controller', { fn: getController, args: ['device-id'] }],
  ['set-IP', { fn: setIP, args: ['device-id', 'ip-address', 'subnet', 'gateway'] }],
  ['get-time', { fn: getTime, args: ['device-id'] }],
  ['set-time', { fn: setTime, args: ['device-id', 'datetime'] }]
])

export function exec (cmd) {
  return cmd.fn()
}

function getAllControllers () {
  return uhppote.GetAllControllers()
}

function getController () {
  const deviceID = document.querySelector('input#device-id').value

  return uhppote.GetController(deviceID)
}

function setIP () {
  const deviceID = document.querySelector('input#device-id').value
  const address = document.querySelector('input#ip-address').value
  const netmask = document.querySelector('input#subnet').value
  const gateway = document.querySelector('input#gateway').value

  return uhppote.SetIP(deviceID, address, netmask, gateway)
}

function getTime () {
  const deviceID = document.querySelector('input#device-id').value

  return uhppote.GetTime(deviceID)
}

function setTime () {
  const deviceID = document.querySelector('input#device-id').value
  const datetime = document.querySelector('input#datetime').value

  if (datetime === '') {
    return uhppote.SetTime(deviceID, new Date())
  } else {
    return uhppote.SetTime(deviceID, new Date(datetime))
  }
}
