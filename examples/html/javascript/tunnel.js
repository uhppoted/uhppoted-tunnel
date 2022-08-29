import * as uhppote from './uhppote.js'

export function initialise () {
  const fields = document.querySelectorAll('input[data-tag]:not([data-tag=""])')

  fields.forEach(f => {
    f.value = get(f.dataset.tag)
  })
}

const vtable = new Map([
  ['get-all-controllers', { fn: getAllControllers }],
  ['get-controller', { fn: getController }],
  ['set-IP', { fn: setIP }],
  ['get-time', { fn: getTime }],
  ['set-time', { fn: setTime }]
])

export function clear () {
  document.querySelector('#request textarea').value = ''
  document.querySelector('#reply textarea').value = ''
  document.querySelector('#response textarea').value = ''

  warn()
}

export function exec (cmd) {
  document.querySelector('#request textarea').value = ''
  document.querySelector('#reply textarea').value = ''
  document.querySelector('#response textarea').value = ''

  warn()

  try {
    const objects = document.querySelector('#response textarea')

    if (vtable.has(cmd)) {
      const f = vtable.get(cmd)

      f.fn().then(response => {
        objects.value = JSON.stringify(response, null, '  ')
      })
    } else {
      warn(`${cmd}: invalid command`)
    }
  } catch (err) {
    warn(err)
  }
}

function getAllControllers () {
  stash([])

  return uhppote.GetAllControllers()
}

function getController () {
  const deviceID = document.querySelector('input#device-id').value

  stash(['device-id'])

  return uhppote.GetController(deviceID)
}

function setIP () {
  const deviceID = document.querySelector('input#device-id').value
  const address = document.querySelector('input#ip-address').value
  const netmask = document.querySelector('input#subnet').value
  const gateway = document.querySelector('input#gateway').value

  stash(['device-id', 'ip-address', 'subnet', 'gateway'])

  return uhppote.SetIP(deviceID, address, netmask, gateway)
}

function getTime () {
  const deviceID = document.querySelector('input#device-id').value

  stash(['device-id'])

  return uhppote.GetTime(deviceID)
}

function setTime () {
  const deviceID = document.querySelector('input#device-id').value
  const datetime = document.querySelector('input#datetime').value

  stash(['device-id'])

  if (datetime === '') {
    return uhppote.SetTime(deviceID, new Date())
  } else {
    return uhppote.SetTime(deviceID, new Date(datetime))
  }
}

function stash (list) {
  const f = function (e) {
    return {
      tag: e.dataset.tag,
      value: e.value
    }
  }

  list.map(id => document.querySelector(`input#${id}`))
    .map(e => f(e))
    .filter(o => o.tag)
    .forEach(o => put(o.tag, o.value))
}

function put (tag, value) {
  localStorage.setItem(tag, JSON.stringify(value))
}

function get (tag) {
  const value = localStorage.getItem(tag)

  return value ? JSON.parse(value) : ''
}

function warn (err) {
  const message = document.getElementById('message')

  if (err) {
    message.innerHTML = err
  } else {
    message.innerHTML = ''
  }
}
