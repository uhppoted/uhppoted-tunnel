import * as main from './main.js'

const USES = new Map([...main.COMMANDS].map(([u, v]) => [u, v.args]))

export function initialise () {
  unstash()
}

export function clear () {
  document.querySelector('#request textarea').value = ''
  document.querySelector('#reply textarea').value = ''
  document.querySelector('#response textarea').value = ''

  warn()
}

export function select (event) {
  const cmd = event.currentTarget.dataset.cmd
  const uses = USES.get(cmd)
  const sections = document.querySelectorAll('div.section')
  const fields = document.querySelectorAll('input[data-tag]:not([data-tag=""])')
  const labels = [...fields].map(e => document.querySelector(`label[for="${e.id}"]`)).filter(l => l !== null)
  const command = document.querySelector('input#cmd')
  const execute = document.querySelector('button#execute')

  // .. show/hide sections
  sections.forEach(section => {
    for (const u of uses) {
      if (section.querySelector(`[data-tag="${u}"]`)) {
        section.classList.add('visible')
        return
      }
    }

    section.classList.remove('visible')
  })

  // .. enable/disable fields
  fields.forEach(e => { e.disabled = true })
  labels.forEach(e => { e.classList.add('disabled') })

  uses.forEach(u => {
    const e = document.querySelector(`[data-tag^="${u}"]`)
    if (e) {
      e.disabled = false

      const label = document.querySelector(`label[for="${e.id}"]`)
      if (label) {
        label.classList.remove('disabled')
      }
    }
  })

  // .. weekdays are a special snowflake
  const label = document.querySelector('label[for="time-profile-weekdays"]')
  if (label) {
    label.classList.add('disabled')

    const weekdays = [...document.querySelectorAll('div#time-profile-weekdays input')].map(e => e.dataset.tag)
    for (const day of weekdays) {
      if (uses.includes(day)) {
        label.classList.remove('disabled')
        break
      }
    }
  }

  // .. set command header

  command.value = cmd
  execute.dataset.command = cmd
  execute.disabled = false
}

export function execute (event) {
  try {
    const cmd = event.currentTarget.dataset.command
    const response = document.querySelector('#response textarea')

    if (cmd) {
      document.querySelector('#request textarea').value = ''
      document.querySelector('#reply textarea').value = ''
      document.querySelector('#response textarea').value = ''

      warn()
      stash(cmd)

      main
        .exec(cmd)
        .then(v => {
          response.value = JSON.stringify(v, null, '  ')
        })
        .catch(err => warn(`${err}`))
    }
  } catch (err) {
    console.error(err)
    warn(err)
  }
}

function stash (cmd) {
  const list = main.COMMANDS.get(cmd).args

  list
    .map(tag => document.querySelector(`[data-tag="${tag}"]`))
    .filter(e => e.dataset.tag !== 'datetime')
    .forEach(e => put(e.dataset.tag, e.type === 'checkbox' ? e.checked : e.value))
}

function unstash () {
  const fields = [...document.querySelectorAll('input[data-tag]:not([data-tag=""])')]

  fields
    .filter(f => f.dataset.tag !== 'datetime')
    .forEach(f => {
      if (f.type === 'checkbox') {
        f.checked = get(f.dataset.tag)
      } else {
        f.value = get(f.dataset.tag)
      }
    })
}

function get (tag) {
  const value = localStorage.getItem(tag)

  return value ? JSON.parse(value) : ''
}

function put (tag, value) {
  localStorage.setItem(tag, JSON.stringify(value))
}

function warn (err) {
  const message = document.getElementById('message')

  if (err) {
    message.innerHTML = err
  } else {
    message.innerHTML = ''
  }
}
