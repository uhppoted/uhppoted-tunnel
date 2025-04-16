import * as commands from './commands.js'

export const COMMANDS = new Map([
  ['get-all-controllers', { fn: commands.getAllControllers, args: [] }],
  ['get-controller', { fn: commands.getController, args: ['controller'] }],
  ['set-IP', { fn: commands.setIP, args: ['controller', 'IP.address', 'IP.netmask', 'IP.gateway'] }],
  ['get-time', { fn: commands.getTime, args: ['controller'] }],
  ['set-time', { fn: commands.setTime, args: ['controller', 'datetime'] }],
  ['get-listener', { fn: commands.getListener, args: ['controller'] }],
  ['set-listener', { fn: commands.setListener, args: ['controller', 'listener', 'interval'] }],
  ['get-door-control', { fn: commands.getDoorControl, args: ['controller', 'door.id'] }],
  ['set-door-control', { fn: commands.setDoorControl, args: ['controller', 'door.id', 'door.mode', 'door.delay'] }],
  ['set-door-passcodes', {
    fn: commands.setDoorPasscodes,
    args: ['controller',
      'door.id',
      'door.passcode1',
      'door.passcode2',
      'door.passcode3',
      'door.passcode4']
  }],
  ['get-status', { fn: commands.getStatus, args: ['controller'] }],
  ['open-door', { fn: commands.openDoor, args: ['controller', 'door.id'] }],
  ['get-cards', { fn: commands.getCards, args: ['controller'] }],
  ['get-card', { fn: commands.getCard, args: ['controller', 'card.number'] }],
  ['get-card-by-index', { fn: commands.getCardByIndex, args: ['controller', 'card.index'] }],
  ['put-card', {
    fn: commands.putCard,
    args: ['controller',
      'card.number',
      'card.start-date',
      'card.end-date',
      'card.doors.1',
      'card.doors.2',
      'card.doors.3',
      'card.doors.4',
      'card.PIN']
  }],

  ['delete-card', { fn: commands.deleteCard, args: ['controller', 'card.number'] }],
  ['delete-all-cards', { fn: commands.deleteAllCards, args: ['controller'] }],
  ['get-event', { fn: commands.getEvent, args: ['controller', 'events.index'] }],
  ['get-event-index', { fn: commands.getEventIndex, args: ['controller'] }],
  ['set-event-index', { fn: commands.setEventIndex, args: ['controller', 'events.index'] }],
  ['record-special-events', { fn: commands.recordSpecialEvents, args: ['controller', 'events.record-special-events'] }],
  ['get-time-profile', { fn: commands.getTimeProfile, args: ['controller', 'time-profile.id'] }],
  ['set-time-profile', {
    fn: commands.setTimeProfile,
    args: ['controller',
      'time-profile.id',
      'time-profile.start-date',
      'time-profile.end-date',
      'time-profile.monday',
      'time-profile.tuesday',
      'time-profile.wednesday',
      'time-profile.thursday',
      'time-profile.friday',
      'time-profile.saturday',
      'time-profile.sunday',
      'time-profile.segment.1.start',
      'time-profile.segment.1.end',
      'time-profile.segment.2.start',
      'time-profile.segment.2.end',
      'time-profile.segment.3.start',
      'time-profile.segment.3.end',
      'time-profile.linked-profile.id']
  }],

  ['delete-all-time-profiles', { fn: commands.deleteAllTimeProfiles, args: ['controller'] }],

  ['add-task', {
    fn: commands.addTask,
    args: ['controller',
      'task.start-date',
      'task.end-date',
      'task.monday',
      'task.tuesday',
      'task.wednesday',
      'task.thursday',
      'task.friday',
      'task.saturday',
      'task.sunday',
      'task.start-time',
      'task.door',
      'task.type',
      'task.more-cards']
  }],
  ['refresh-tasklist', { fn: commands.refreshTaskList, args: ['controller'] }],
  ['clear-tasklist', { fn: commands.clearTaskList, args: ['controller'] }],
  ['set-pc-control', { fn: commands.setPCControl, args: ['controller', 'pc-control'] }],
  ['set-interlock', { fn: commands.setInterlock, args: ['controller', 'interlock'] }],
  ['activate-keypads', { fn: commands.activateKeypads, args: ['controller', 'reader1', 'reader2', 'reader3', 'reader4'] }],
  ['get-antipassback', { fn: commands.getAntiPassback, args: ['controller'] }],
  ['set-antipassback', { fn: commands.setAntiPassback, args: ['controller', 'antipassback'] }],
  ['restore-default-parameters', { fn: commands.restoreDefaultParameters, args: ['controller'] }]
])

export function exec (cmd) {
  if (COMMANDS.has(cmd)) {
    const c = COMMANDS.get(cmd)

    return commands.exec(c.fn, ...c.args)
  } else {
    throw new Error(`invalid command '${cmd}'`)
  }
}
