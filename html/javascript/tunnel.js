import * as encode from './encode.js'

export function initialise() {
}

export function exec(cmd) {
  warn()
  switch (cmd) {
    case 'get-devices':
       post(encode.get_devices())
       break

    default:
      warn(`${cmd}: invalid command`)
  }
}

function post(rq) {
  const hex = bin2hex(rq)
  const text = document.querySelector('#request textarea')

  text.value = hex
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

  for (let i=0; i < chunks.length; i += 2) {
    lines.push([chunks[i], chunks[i+1]].join('  '))
  }
  
  return lines.join(`\n`)

  // const f = function* chunks(array,N) {
  //    for (let i=0; i < array.length; i += N) {
  //        yield array.slice(i, i + N);
  //    }
  // }
}
