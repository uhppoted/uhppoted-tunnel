import * as uhppote from './uhppote.js'

export function exec (fn, ...args) {
  return fn(...args)
}

export function getAllControllers () {
  return uhppote.GetAllControllers()
}

export function getController (controller) {
  controller = document.querySelector(`input#${controller}`).value

  return uhppote.GetController(controller)
}

export function setIP (controller, address, netmask, gateway) {
  controller = document.querySelector(`input#${controller}`).value
  address = document.querySelector(`input#${address}`).value
  netmask = document.querySelector(`input#${netmask}`).value
  gateway = document.querySelector(`input#${gateway}`).value

  return uhppote.SetIP(controller, address, netmask, gateway)
}

export function getTime (controller) {
  controller = document.querySelector(`input#${controller}`).value

  return uhppote.GetTime(controller)
}

export function setTime (controller, datetime) {
  controller = document.querySelector(`input#${controller}`).value
  const dt = document.querySelector(`input#${datetime}`).value

  if (dt === '') {
    return uhppote.SetTime(controller, new Date())
  } else {
    return uhppote.SetTime(controller, new Date(datetime))
  }
}

// export function getStatus(controller) {
//  controller = document.querySelector(`input#${controller}`).value
//
//  return uhppote.GetStatus(controller)
// }
//
// export function getListener(controller) {
//  controller = document.querySelector(`input#${controller}`).value
//
//  return uhppote.GetListener(controller)
// }
//
// export function setListener(controller, addr) {
//  controller = document.querySelector(`input#${controller}`).value
//  addr = document.querySelector(`input#${addr}`).value
//
//  const address = addr.match(/^(.*?):([0-9])+$/)[1]
//  const port = addr.match(/^(.*?):([0-9])+$/)[2]
//
//  return uhppote.SetListener(controller, address, port)
// }
//
// export function getDoorControl(controller, door) {
//  controller = document.querySelector(`input#${controller}`).value
//  door = document.querySelector(`input#${door}`).value
//
//  return uhppote.GetDoorControl(controller,door)
// }
//
// export function setDoorControl(controller, door, mode, delay) {
//  controller = document.querySelector(`input#${controller}`).value
//  door = document.querySelector(`input#${door}`).value
//  mode = document.querySelector(`input#${mode}`).value
//  delay = document.querySelector(`input#${delay}`).value
//
//  return uhppote.SetDoorControl(controller, door, mode, delay)
// }
//
// export function openDoor(controller, door) {
//  controller = document.querySelector(`input#${controller}`).value
//  door = document.querySelector(`input#${door}`).value
//
//  return uhppote.OpenDoor(controller,door)
// }
//
// export function getCards(controller) {
//  controller = document.querySelector(`input#${controller}`).value
//
//  return uhppote.GetCards(controller)
// }
//
// export function getCard(controller, card) {
//  controller = document.querySelector(`input#${controller}`).value
//  card = document.querySelector(`input#${card}`).value
//
//  response = uhppote.GetCard(controller, card)
//  if (response.cardNumber === 0) {
//      throw new Error(`card ${card} not found`)
//  }
//
//  return response
// }
//
// export function getCardByIndex(controller, index) {
//  controller = document.querySelector(`input#${controller}`).value
//  index = document.querySelector(`input#${index}`).value
//
//  response = uhppote.GetCardByIndex(controller,index)
//  if (response.cardNumber === 0) {
//      throw new Error(`card @ index ${index} not found`)
//  } else if (response.cardNumber === 0xffffffff) {
//      throw new Error(`card @ index ${index} deleted`)
//  }
//
//  return response
// }
//
// export function putCard(controller, card, start, end, door1, door2, door3, door4) {
//  controller = document.querySelector(`input#${controller}`).value
//  card = document.querySelector(`input#${card}`).value
//  start = document.querySelector(`input#${start}`).value
//  end = document.querySelector(`input#${end}`).value
//  door1 = document.querySelector(`input#${door1}`).value
//  door2 = document.querySelector(`input#${door2}`).value
//  door3 = document.querySelector(`input#${door3}`).value
//  door4 = document.querySelector(`input#${door4}`).value
//
//  return uhppote.PutCard(controller, card, new Date(start), new Date(end), door1, door2, door3, door4)
// }
//
// export function deleteCard(controller, card) {
//  controller = document.querySelector(`input#${controller}`).value
//  card = document.querySelector(`input#${card}`).value
//
//  return uhppote.DeleteCard(controller, card)
// }
//
// export function deleteAllCards(controller) {
//  controller = document.querySelector(`input#${controller}`).value
//
//  return uhppote.DeleteAllCards(controller)
// }
//
// export function getEvent(controller, index) {
//  controller = document.querySelector(`input#${controller}`).value
//  index = document.querySelector(`input#${index}`).value
//
//  response = uhppote.GetEvent(controller,index)
//  if (response.eventType === 0xff) {
//      throw new Error(`event @ index ${index} overwritten`)
//  } else if (response.index === 0) {
//      throw new Error(`event @ index ${index} not found`)
//  }
//
//  return response
// }
//
// export function getEventIndex(controller) {
//  controller = document.querySelector(`input#${controller}`).value
//
//  return uhppote.GetEventIndex(controller)
// }
//
// export function setEventIndex(controller, index) {
//  controller = document.querySelector(`input#${controller}`).value
//  index = document.querySelector(`input#${index}`).value
//
//  return uhppote.SetEventIndex(controller, index)
// }
//
// export function recordSpecialEvents(controller, enabled) {
//  controller = document.querySelector(`input#${controller}`).value
//  enabled = document.querySelector(`input#${enabled}`).checked
//
//  return uhppote.RecordSpecialEvents(controller, enabled)
// }
//
// export function getTimeProfile(controller, profileID) {
//  controller = document.querySelector(`input#${controller}`).value
//  profileID = document.querySelector(`input#${profileID}`).value
//
//  response = uhppote.GetTimeProfile(controller, profileID)
//  if (response.profileId === 0) {
//      throw new Error(`time profile ${profileID} not defined`)
//  }
//
//  return response
// }
//
/// / export function setTimeProfile(controller,
//   profileID,
//   start, end,
//   monday, tuesday, wednesday, thursday, friday, saturday, sunday,
//   segment1start, segment1end,
//   segment2start, segment2end,
//   segment3start, segment3end,
//   linkedProfileID) {
//   controller = document.querySelector(`input#${controller}`).value
//   profileID = document.querySelector(`input#${profileID}`).value
//   start = document.querySelector(`input#${start}`).value
//   end = document.querySelector(`input#${end}`).value
//
//   monday = document.querySelector(`input#${monday}`).checked
//   tuesday = document.querySelector(`input#${tuesday}`).checked
//   wednesday = document.querySelector(`input#${wednesday}`).checked
//   thursday = document.querySelector(`input#${thursday}`).checked
//   friday = document.querySelector(`input#${friday}`).checked
//   saturday = document.querySelector(`input#${saturday}`).checked
//   sunday = document.querySelector(`input#${sunday}`).checked
//
//   segment1start = document.querySelector(`input#${segment1start}`).value
//   segment1end = document.querySelector(`input#${segment1end}`).value
//   segment2start = document.querySelector(`input#${segment2start}`).value
//   segment2end = document.querySelector(`input#${segment2end}`).value
//   segment3start = document.querySelector(`input#${segment3start}`).value
//   segment3end = document.querySelector(`input#${segment3end}`).value
//
//   linkedProfileID = document.querySelector(`input#${linkedProfileID}`).value
//
//   return uhppote.SetTimeProfile(controller,
//         profileID,
//         new Date(start), new Date(end),
//         monday, tuesday, wednesday, thursday, friday, saturday, sunday,
//         uhppote.HHmm(segment1start), uhppote.HHmm(segment1end),
//         uhppote.HHmm(segment2start), uhppote.HHmm(segment2end),
//         uhppote.HHmm(segment3start), uhppote.HHmm(segment3end),
//         linkedProfileID)
// }
//
// export function deleteAllTimeProfiles(controller) {
//   controller = document.querySelector(`input#${controller}`).value
//
//   return uhppote.DeleteAllTimeProfiles(controller)
// }
//
// export function addTask(controller,
//   start, end,
//   monday, tuesday, wednesday, thursday, friday, saturday, sunday,
//   time,
//   door, taskType, moreCards) {
//   controller = document.querySelector(`input#${controller}`).value
//   start = document.querySelector(`input#${start}`).value
//   end = document.querySelector(`input#${end}`).value
//
//   monday = document.querySelector(`input#${monday}`).checked
//   tuesday = document.querySelector(`input#${tuesday}`).checked
//   wednesday = document.querySelector(`input#${wednesday}`).checked
//   thursday = document.querySelector(`input#${thursday}`).checked
//   friday = document.querySelector(`input#${friday}`).checked
//   saturday = document.querySelector(`input#${saturday}`).checked
//   sunday = document.querySelector(`input#${sunday}`).checked
//
//   time = document.querySelector(`input#${time}`).value
//   door = document.querySelector(`input#${door}`).value
//   taskType = document.querySelector(`input#${taskType}`).value
//   moreCards = document.querySelector(`input#${moreCards}`).value
//
//   return uhppote.AddTask(controller,
//         Date(start), Date(end),
//         monday, tuesday, wednesday, thursday, friday, saturday, sunday,
//         time,
//         door,
//         taskType,
//         moreCards)
// }
//
// export function refreshTaskList(controller) {
//   controller = document.querySelector(`input#${controller}`).value
//
//   return uhppote.RefreshTasklist(controller)
// }
//
// export function clearTaskList(controller) {
//   controller = document.querySelector(`input#${controller}`).value
//
//   return uhppote.ClearTasklist(controller)
// }
