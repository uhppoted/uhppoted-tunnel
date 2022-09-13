let REQUESTID = 0

export function broadcast (bytes) {
  debug([bytes], '#request textarea')

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
          return Promise.reject(response.statusText)
      }
    })
    .then(reply => {
      debug(reply.replies, '#reply textarea')
      return reply.replies
    })
    .catch(err => {
      throw new Error(err)
    })
}

export function send (bytes, nowait) {
  debug([bytes], '#request textarea')

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
          return Promise.reject(response.statusText)
      }
    })
    .then(reply => {
      if (reply.reply) {
        debug([reply.reply], '#reply textarea')
        return reply.reply
      }

      return null
    })
    .catch(err => {
      throw new Error(err)
    })
}

function nextID () {
  REQUESTID++

  return REQUESTID
}

function debug (messages, selector) {
  const hex = messages.map(m => bin2hex(m)).join('\n\n')

  const textarea = document.querySelector(selector)
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
}
