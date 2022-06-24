import * as encode from './encode.js'
import * as decode from './decode.js'

export function initialise() {
}

var REQUESTID = 0

export function exec(cmd) {
  document.querySelector('#request textarea').value = ''
  document.querySelector('#reply textarea').value = ''
  document.querySelector('#response textarea').value = ''

  warn()
  
  switch (cmd) {
    case 'get-devices':
       post(encode.GetDevices())
       break

    default:
      warn(`${cmd}: invalid command`)
  }
}

function post(bytes) {
  const hex = bin2hex(bytes)
  const debug = document.querySelector('#request textarea')

  debug.value = hex

  const rq = {
    ID: nextID(),
    request: [...bytes],
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
            break

          default:
            response.text().then(w => { 
              warn(new Error(w)) 
            })
        }
    })
    .then(reply => {
      result(reply.reply)
    })
    .catch(function (err) {
      warn(`${err.message.toLowerCase()}`)
    })
    .finally(() => {
    })
}

function result(bytes) {
  const hex = bin2hex(bytes)
  const debug = document.querySelector('#reply textarea')
  const object = document.querySelector('#response textarea')

  const o = decode.GetDevice(bytes)

  debug.value = hex
  object.value = JSON.stringify(o, null, '  ')
}

function nextID() {
  REQUESTID++

  return REQUESTID
}

function warn(err) {
   const message = document.getElementById('message')

   if (err) {
      message.innerHTML = err
   } else {
      message.innerHTML = ''
   }
}

function bin2hex(bytes) { 
  const chunks = [...bytes]
                 .map(x => x.toString(16).padStart(2, '0'))
                 .join('')
                 .match(/.{1,16}/g)
                 .map(l => l.match(/.{1,2}/g).join(' '))

  const lines = []
  while (chunks.length > 0) {
    lines.push(chunks.splice(0,2).join('  '))
  }

  return lines.join(`\n`)

  // const f = function* chunks(array,N) {
  //    for (let i=0; i < array.length; i += N) {
  //        yield array.slice(i, i + N);
  //    }
  // }
}
