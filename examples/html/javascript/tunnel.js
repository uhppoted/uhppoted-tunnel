import * as encode from './encode.js'
import * as decode from './decode.js'

export function initialise () {
}

let REQUESTID = 0

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

  switch (cmd) {
    case 'get-devices':
      post(encode.GetDevices, '500ms')
      break

    case 'get-device':
      post(encode.GetDevice, '0s', 'device-id')
      break

    case 'set-address':
      post(encode.SetIP, '0.1ms', 'device-id', 'ip-address', 'subnet', 'gateway')
      break

    default:
      warn(`${cmd}: invalid command`)
  }
}

function post (f, timeout, ...args) {
  const bytes = f(...args.map(a => document.querySelector(`input#${a}`).value))
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
    responses.push(decode.GetDevice(reply))
  }

  debug.value = hex.join('\n\n')
  objects.value = JSON.stringify(responses, null, '  ')
}

function nextID () {
  REQUESTID++

  return REQUESTID
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
