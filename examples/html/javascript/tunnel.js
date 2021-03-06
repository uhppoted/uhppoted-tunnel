import * as encoder from './encode.js'
import * as decoder from './decode.js'

export function initialise () {
  const fields = document.querySelectorAll('input[data-tag]:not([data-tag=""])')

  fields.forEach(f => {
    f.value = get(f.dataset.tag)
  })
}

let REQUESTID = 0

/* eslint-disable */
const vtable = new Map([
  ['get-devices', { encode: encoder.GetDevices, args: [],                                               timeout: '500ms' }],
  ['get-device',  { encode: encoder.GetDevice,  args: ['device-id'],                                    timeout: '0s'    }],
  ['set-address', { encode: encoder.SetIP,      args: ['device-id', 'ip-address', 'subnet', 'gateway'], timeout: '0.1ms' }],
  ['get-time',    { encode: encoder.GetTime,    args: ['device-id'],                                    timeout: '0s'    }],
  ['set-time',    { encode: encoder.SetTime,    args: ['device-id', 'time'],                            timeout: '0s'    }],
])
/* eslint-enable */

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
    if (vtable.has(cmd)) {
      const f = vtable.get(cmd)
      const bytes = f.encode(...f.args.map(a => document.querySelector(`input#${a}`).value))

      stash(f.args)
      post(bytes, f.timeout)
    } else {
      warn(`${cmd}: invalid command`)
    }
  } catch (err) {
    warn(err)
  }
}

function post (bytes, timeout) {
  const hex = bin2hex(bytes)
  const debug = document.querySelector('#request textarea')

  debug.value = hex

  const rq = {
    ID: nextID(),
    wait: timeout,
    request: [...bytes]
  }

  const request = {
    method: 'POST',
    mode: 'cors',
    cache: 'no-cache',
    credentials: 'same-origin',
    headers: { 'Content-Type': 'application/json' },
    redirect: 'follow',
    referrerPolicy: 'no-referrer',
    body: JSON.stringify(rq)
  }

  fetch('/udp', request)
    .then(response => {
      switch (response.status) {
        case 200:
          return response.json()

        default:
          response.text().then(w => {
            warn(new Error(w))
          })
      }
    })
    .then(reply => {
      result(reply.replies)
    })
    .catch(function (err) {
      warn(`${err.message.toLowerCase()}`)
    })
    .finally(() => {
    })
}

function result (replies) {
  const debug = document.querySelector('#reply textarea')
  const objects = document.querySelector('#response textarea')
  const hex = []
  const responses = []

  for (const reply of replies) {
    hex.push(bin2hex(reply))
    responses.push(decoder.decode(reply))
  }

  debug.value = hex.join('\n\n')
  objects.value = JSON.stringify(responses, null, '  ')
}

function nextID () {
  REQUESTID++

  return REQUESTID
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

function bin2hex (bytes) {
  const chunks = [...bytes]
    .map(x => x.toString(16).padStart(2, '0'))
    .join('')
    .match(/.{1,16}/g)
    .map(l => l.match(/.{1,2}/g).join(' '))

  const lines = []
  while (chunks.length > 0) {
    lines.push(chunks.splice(0, 2).join('  '))
  }

  return lines.join('\n')

  // const f = function* chunks(array,N) {
  //    for (let i=0; i < array.length; i += N) {
  //        yield array.slice(i, i + N);
  //    }
  // }
}
