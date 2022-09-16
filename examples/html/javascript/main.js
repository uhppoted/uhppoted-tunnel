import * as commands from './commands.js'

const COMMANDS = new Map([
  ['get-all-controllers', { fn: commands.getAllControllers, args: [] }],
  ['get-controller', { fn: commands.getController, args: ['controller'] }],
  ['set-IP', { fn: commands.setIP, args: ['controller', 'ip-address', 'subnet', 'gateway'] }],
  ['get-time', { fn: commands.getTime, args: ['controller'] }],
  ['set-time', { fn: commands.setTime, args: ['controller', 'datetime'] }],
  ['get-listener', { fn: commands.getListener, args: ['controller'] }],
  ['set-listener', { fn: commands.setListener, args: ['controller', 'listener'] }],
  ['get-door-control', { fn: commands.getDoorControl, args: ['controller', 'door'] }],
  ['set-door-control', { fn: commands.setDoorControl, args: ['controller', 'door', 'mode', 'delay'] }],
  ['get-status', { fn: commands.getStatus, args: ['controller'] }],
  ['open-door', { fn: commands.openDoor, args: ['controller', 'door'] }],
  ['get-cards', { fn: commands.getCards, args: ['controller'] }],
  ['get-card', { fn: commands.getCard, args: ['controller', 'card-number'] }],
  ['get-card-by-index', { fn: commands.getCardByIndex, args: ['controller', 'card-index'] }],
  ['put-card', { fn: commands.putCard, args: ['controller', 'card-number', 'card-start-date', 'card-end-date', 'card-door-1', 'card-door-2', 'card-door-3', 'card-door-4'] }],
  ['delete-card', { fn: commands.deleteCard, args: ['controller', 'card-number'] }],
  ['delete-all-cards', { fn: commands.deleteAllCards, args: ['controller'] }],
  ['get-event', { fn: commands.getEvent, args: ['controller', 'event-index'] }],
  ['get-event-index', { fn: commands.getEventIndex, args: ['controller'] }],
  ['set-event-index', { fn: commands.setEventIndex, args: ['controller', 'event-index'] }],
  ['record-special-events', { fn: commands.recordSpecialEvents, args: ['controller', 'record-special-events'] }]
  // ['get-timeprofile', { fn: commands.getTimeProfile, args: ['controller', 'profileID'] }],
  // ['set-timeprofile', {
  //   fn: commands.setTimeProfile,
  //   args: ['controller',
  //     'profileID',
  //     'start', 'end',
  //     'monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday',
  //     'segment1start', 'segment1end',
  //     'segment2start', 'segment2end',
  //     'segment3start', 'segment3end',
  //     'linkedProfileID']
  // }],
  // ['delete-all-time-profiles', { fn: commands.deleteAllTimeProfiles, args: ['controller'] }],
  // ['add-task', {
  //   fn: commands.addTask,
  //   args: ['controller',
  //     'start', 'end',
  //     'monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday',
  //     'time',
  //     'door', 'taskType', 'moreCards']
  // }],
  // ['refresh-tasklist', { fn: commands.refreshTaskList, args: ['controller'] }],
  // ['clear-tasklist', { fn: commands.clearTaskList, args: ['controller'] }]
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
        .catch(err => {
          warn(`${err}`)
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
