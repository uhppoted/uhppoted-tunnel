let REQUESTID = 0

export function broadcast (bytes) {
  debug(bytes)

  const rq = {
    ID: nextID(),
    wait: '500ms',
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

  return fetch('/udp/broadcast', request)
    .then(response => {
      switch (response.status) {
        case 200:
          return response.json()

        default:
          response.text().then(w => {
            throw new Error(w)
          })
      }
    })
    .then(reply => {
      debug2(reply.replies)
      return reply.replies
    })
}

export function send (bytes, nowait) {
  debug(bytes)

  const rq = {
    ID: nextID(),
    wait: !nowait,
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

  return fetch('/udp/send', request)
    .then(response => {
      switch (response.status) {
        case 200:
          return response.json()

        default:
          response.text().then(w => {
            throw new Error(w)
          })
      }
    })
    .then(reply => {
      if (reply.reply) {
        debug2([reply.reply])
        return reply.reply
      }

      return null
    })
}

function nextID () {
  REQUESTID++

  return REQUESTID
}

function debug (request) {
  const textarea = document.querySelector('#request textarea')
  if (textarea) {
    textarea.value = bin2hex(request)
  }
}

function debug2 (replies) {
  const hex = replies.map(r => bin2hex(r)).join('\n\n')
  const textarea = document.querySelector('#reply textarea')
  if (textarea) {
    textarea.value = hex
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
