import * as uhppote from './uhppote.js'

export function exec (fn, ...args) {
  return fn(...args)
}

export function getAllControllers () {
  return uhppote.GetAllControllers()
}

export function getController (controller) {
  controller = arg(controller)

  return uhppote.GetController(controller)
}

export function setIP (controller, address, netmask, gateway) {
  controller = arg(controller)
  address = arg(address)
  netmask = arg(netmask)
  gateway = arg(gateway)

  return uhppote.SetIP(controller, address, netmask, gateway)
}

export function getTime (controller) {
  controller = arg(controller)

  return uhppote.GetTime(controller)
}

export function setTime (controller, datetime) {
  controller = arg(controller)
  const dt = arg(datetime)

  if (dt === '') {
    return uhppote.SetTime(controller, new Date())
  } else {
    return uhppote.SetTime(controller, new Date(dt))
  }
}

export function getListener (controller) {
  controller = arg(controller)

  return uhppote.GetListener(controller)
}

export function setListener (controller, listener, interval) {
  controller = arg(controller)
  listener = arg(listener)
  interval = arg(interval)

  const address = listener.match(/^(.*?):([0-9]+)$/)[1]
  const port = listener.match(/^(.*?):([0-9]+)$/)[2]

  return uhppote.SetListener(controller, address, port, interval)
}

export function getDoorControl (controller, door) {
  controller = arg(controller)
  door = arg(door)

  return uhppote.GetDoorControl(controller, door)
}

export function setDoorControl (controller, door, mode, delay) {
  controller = arg(controller)
  door = arg(door)
  mode = arg(mode)
  delay = arg(delay)

  return uhppote.SetDoorControl(controller, door, mode, delay)
}

export function openDoor (controller, door) {
  controller = arg(controller)
  door = arg(door)

  return uhppote.OpenDoor(controller, door)
}

export function getStatus (controller) {
  controller = arg(controller)

  return uhppote.GetStatus(controller)
}

export function getCards (controller) {
  controller = arg(controller)

  return uhppote.GetCards(controller)
}

export function getCard (controller, card) {
  controller = arg(controller)
  card = arg(card)

  return uhppote.GetCard(controller, card).then(response => {
    if (response.cardNumber === 0) {
      throw new Error(`card ${card} not found`)
    }

    return response
  })
}

export function getCardByIndex (controller, index) {
  controller = arg(controller)
  index = arg(index)

  return uhppote.GetCardByIndex(controller, index).then(response => {
    if (response.cardNumber === 0) {
      throw new Error(`card @ index ${index} not found`)
    } else if (response.cardNumber === 0xffffffff) {
      throw new Error(`card @ index ${index} deleted`)
    }

    return response
  })
}

export function putCard (controller, card, start, end, door1, door2, door3, door4, PIN) {
  controller = arg(controller)
  card = arg(card)
  start = arg(start)
  end = arg(end)
  door1 = arg(door1)
  door2 = arg(door2)
  door3 = arg(door3)
  door4 = arg(door4)
  PIN = arg(PIN)

  return uhppote.PutCard(controller, card, start, end, door1, door2, door3, door4, PIN)
}

export function deleteCard (controller, card) {
  controller = arg(controller)
  card = arg(card)

  return uhppote.DeleteCard(controller, card)
}

export function deleteAllCards (controller) {
  controller = arg(controller)

  return uhppote.DeleteAllCards(controller)
}

export function getEvent (controller, index) {
  controller = arg(controller)
  index = arg(index)

  return uhppote.GetEvent(controller, index).then(response => {
    if (response.eventType === 0xff) {
      throw new Error(`event @ index ${index} overwritten`)
    } else if (response.index === 0) {
      throw new Error(`event @ index ${index} not found`)
    }

    return response
  })
}

export function getEventIndex (controller) {
  controller = arg(controller)

  return uhppote.GetEventIndex(controller)
}

export function setEventIndex (controller, index) {
  controller = arg(controller)
  index = arg(index)

  return uhppote.SetEventIndex(controller, index)
}

export function recordSpecialEvents (controller, enabled) {
  controller = arg(controller)
  enabled = arg(enabled)

  return uhppote.RecordSpecialEvents(controller, enabled)
}

export function getTimeProfile (controller, profileID) {
  controller = arg(controller)
  profileID = arg(profileID)

  return uhppote.GetTimeProfile(controller, profileID).then(response => {
    if (response.profileId === 0) {
      throw new Error(`time profile ${profileID} not defined`)
    }

    return response
  })
}

export function setTimeProfile (controller,
  profileID,
  start, end,
  monday, tuesday, wednesday, thursday, friday, saturday, sunday,
  segment1start, segment1end,
  segment2start, segment2end,
  segment3start, segment3end,
  linkedProfileID) {
  controller = arg(controller)
  profileID = arg(profileID)
  start = arg(start)
  end = arg(end)

  monday = arg(monday)
  tuesday = arg(tuesday)
  wednesday = arg(wednesday)
  thursday = arg(thursday)
  friday = arg(friday)
  saturday = arg(saturday)
  sunday = arg(sunday)

  segment1start = arg(segment1start)
  segment1end = arg(segment1end)
  segment2start = arg(segment2start)
  segment2end = arg(segment2end)
  segment3start = arg(segment3start)
  segment3end = arg(segment3end)

  linkedProfileID = arg(linkedProfileID)

  return uhppote.SetTimeProfile(controller,
    profileID,
    start, end,
    monday, tuesday, wednesday, thursday, friday, saturday, sunday,
    segment1start, segment1end,
    segment2start, segment2end,
    segment3start, segment3end,
    linkedProfileID)
}

export function deleteAllTimeProfiles (controller) {
  controller = arg(controller)

  return uhppote.DeleteAllTimeProfiles(controller)
}

export function addTask (controller,
  startdate, enddate,
  monday, tuesday, wednesday, thursday, friday, saturday, sunday,
  starttime,
  door, taskType, moreCards) {
  controller = arg(controller)
  startdate = arg(startdate)
  enddate = arg(enddate)

  monday = arg(monday)
  tuesday = arg(tuesday)
  wednesday = arg(wednesday)
  thursday = arg(thursday)
  friday = arg(friday)
  saturday = arg(saturday)
  sunday = arg(sunday)

  starttime = arg(starttime)
  door = arg(door)
  taskType = arg(taskType)
  moreCards = arg(moreCards)

  return uhppote.AddTask(controller,
    startdate, enddate,
    monday, tuesday, wednesday, thursday, friday, saturday, sunday,
    starttime,
    door,
    taskType,
    moreCards)
}

export function refreshTaskList (controller) {
  controller = arg(controller)

  return uhppote.RefreshTasklist(controller)
}

export function clearTaskList (controller) {
  controller = arg(controller)

  return uhppote.ClearTasklist(controller)
}

export function setPCControl (controller, enabled) {
  controller = arg(controller)
  enabled = arg(enabled)

  return uhppote.SetPCControl(controller, enabled)
}

export function setInterlock (controller, interlock) {
  controller = arg(controller)
  interlock = arg(interlock)

  return uhppote.SetInterlock(controller, interlock)
}

export function activateKeypads (controller, reader1, reader2, reader3, reader4) {
  controller = arg(controller)
  reader1 = arg(reader1)
  reader2 = arg(reader2)
  reader3 = arg(reader3)
  reader4 = arg(reader4)

  return uhppote.ActivateKeypads(controller, reader1, reader2, reader3, reader4)
}

export function setDoorPasscodes (controller, door, passcode1, passcode2, passcode3, passcode4) {
  controller = arg(controller)
  door = arg(door)
  passcode1 = arg(passcode1)
  passcode2 = arg(passcode2)
  passcode3 = arg(passcode3)
  passcode4 = arg(passcode4)

  return uhppote.SetDoorPasscodes(controller, door, passcode1, passcode2, passcode3, passcode4)
}

export function getAntiPassback (controller) {
  controller = arg(controller)

  return uhppote.GetAntipassback(controller)
}

export function setAntiPassback (controller, antipassback) {
  controller = arg(controller)
  antipassback = arg(antipassback)

  return uhppote.SetAntipassback(controller, antipassback)
}

export function restoreDefaultParameters (controller) {
  controller = arg(controller)

  return uhppote.RestoreDefaultParameters(controller)
}

function arg (tag) {
  let e = document.querySelector(`[data-tag="${tag}"]`)

  if (e) {
    if (e.type === 'checkbox') {
      return e.checked
    } else {
      return e.value
    }
  }

  e = document.querySelector(`#${tag}`)

  if (e) {
    if (e.type === 'checkbox') {
      return e.checked
    } else {
      return e.value
    }
  }
}
