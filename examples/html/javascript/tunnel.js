import * as commands from './commands.js'

const COMMANDS = new Map([
  ['get-all-controllers', { fn: commands.getAllControllers, args: [] }],
  ['get-controller', { fn: commands.getController, args: ['device-id'] }],
  ['set-IP', { fn: commands.setIP, args: ['device-id', 'ip-address', 'subnet', 'gateway'] }],
  ['get-time', { fn: commands.getTime, args: ['device-id'] }],
  ['set-time', { fn: commands.setTime, args: ['device-id', 'datetime'] }]
])

export function initialise () {
  const fields = document.querySelectorAll('input[data-tag]:not([data-tag=""])')

  fields.forEach(f => {
    f.value = get(f.dataset.tag)
  })
}

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

    if (COMMANDS.has(cmd)) {
      const c = COMMANDS.get(cmd)

      stash(c.args)

      commands.exec(c.fn, ...c.args).then(response => {
        objects.value = JSON.stringify(response, null, '  ')
      })
    } else {
      throw new Error(`invalid command '${cmd}'`)
    }
  } catch (err) {
    warn(err)
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
    .filter(e => e !== null && e.dataset.tag)
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
