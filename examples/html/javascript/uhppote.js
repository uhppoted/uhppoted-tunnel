import * as encode from './encode.js'
import * as decode from './decode.js'
import * as udp from './udp.js'

export function GetAllControllers () {
  const request = encode.GetControllerRequest(0)

  return udp.broadcast(request)
    .then(replies => {
      const list = []

      for (const reply of replies) {
        list.push(decode.GetControllerResponse(reply))
      }

      return list
    })
}

export function GetController (deviceId) {
  const request = encode.GetControllerRequest(deviceId)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetControllerResponse(reply) : null
    })
}

export function SetIP (deviceId, address, netmask, gateway) {
  const request = encode.SetIPRequest(deviceId, address, netmask, gateway)

  return udp.send(request, true)
    .then(() => {
      return true
    })
}

export function GetTime (deviceId) {
  const request = encode.GetTimeRequest(deviceId)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetTimeResponse(reply) : null
    })
}

export function SetTime (deviceId, time) {
  const request = encode.SetTimeRequest(deviceId, time)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetTimeResponse(reply) : null
    })
}

export function GetStatus (deviceId) {
  const request = encode.GetStatusRequest(deviceId)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetStatusResponse(reply) : null
    })
}

export function GetListener (deviceId) {
  const request = encode.GetListenerRequest(deviceId)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetListenerResponse(reply) : null
    })
}

export function SetListener (deviceId, address, port) {
  const request = encode.SetListenerRequest(deviceId, address, port)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetListenerResponse(reply) : null
    })
}

export function GetDoorControl (deviceId, door) {
  const request = encode.GetDoorControlRequest(deviceId, door)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetDoorControlResponse(reply) : null
    })
}

export function SetDoorControl (deviceId, door, mode, delay) {
  const request = encode.SetDoorControlRequest(deviceId, door, mode, delay)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetDoorControlResponse(reply) : null
    })
}

export function OpenDoor (deviceId, door) {
  const request = encode.OpenDoorRequest(deviceId, door)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.OpenDoorResponse(reply) : null
    })
}

export function GetCards (deviceId) {
  const request = encode.GetCardsRequest(deviceId)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetCardsResponse(reply) : null
    })
}

export function GetCard (deviceId, cardNumber) {
  const request = encode.GetCardRequest(deviceId, cardNumber)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetCardResponse(reply) : null
    })
}

export function GetCardByIndex (deviceId, cardIndex) {
  const request = encode.GetCardByIndexRequest(deviceId, cardIndex)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetCardByIndexResponse(reply) : null
    })
}

export function PutCard (deviceId, cardNumber, startDate, endDate, door1, door2, door3, door4, PIN) {
  const request = encode.PutCardRequest(deviceId, cardNumber, startDate, endDate, door1, door2, door3, door4, PIN)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.PutCardResponse(reply) : null
    })
}

export function DeleteCard (deviceId, cardNumber) {
  const request = encode.DeleteCardRequest(deviceId, cardNumber)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.DeleteCardResponse(reply) : null
    })
}

export function DeleteAllCards (deviceId) {
  const request = encode.DeleteCardsRequest(deviceId)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.DeleteAllCardsResponse(reply) : null
    })
}

export function GetEvent (deviceId, eventIndex) {
  const request = encode.GetEventRequest(deviceId, eventIndex)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetEventResponse(reply) : null
    })
}

export function GetEventIndex (deviceId) {
  const request = encode.GetEventIndexRequest(deviceId)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetEventIndexResponse(reply) : null
    })
}

export function SetEventIndex (deviceId, eventIndex) {
  const request = encode.SetEventIndexRequest(deviceId, eventIndex)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetEventIndexResponse(reply) : null
    })
}

export function RecordSpecialEvents (deviceId, enable) {
  const request = encode.RecordSpecialEventsRequest(deviceId, enable)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.RecordSpecialEventsResponse(reply) : null
    })
}

export function GetTimeProfile (deviceId, profileId) {
  const request = encode.GetTimeProfileRequest(deviceId, profileId)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetTimeProfileResponse(reply) : null
    })
}

export function SetTimeProfile (deviceId, profileId, startDate, endDate, monday, tuesday, wednesday, thursday, friday, saturday, sunday, segment1Start, segment1End, segment2Start, segment2End, segment3Start, segment3End, linkedProfileId) {
  const request = encode.SetTimeProfileRequest(deviceId, profileId, startDate, endDate, monday, tuesday, wednesday, thursday, friday, saturday, sunday, segment1Start, segment1End, segment2Start, segment2End, segment3Start, segment3End, linkedProfileId)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetTimeProfileResponse(reply) : null
    })
}

export function DeleteAllTimeProfiles (deviceId) {
  const request = encode.DeleteAllTimeProfilesRequest(deviceId)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.DeleteAllTimeProfilesResponse(reply) : null
    })
}

export function AddTask (deviceId, startDate, endDate, monday, tuesday, wednesday, thursday, friday, saturday, sunday, startTime, door, taskType, moreCards) {
  const request = encode.AddTaskRequest(deviceId, startDate, endDate, monday, tuesday, wednesday, thursday, friday, saturday, sunday, startTime, door, taskType, moreCards)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.AddTaskResponse(reply) : null
    })
}

export function RefreshTasklist (deviceId) {
  const request = encode.RefreshTasklistRequest(deviceId)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.RefreshTasklistResponse(reply) : null
    })
}

export function ClearTasklist (deviceId) {
  const request = encode.ClearTasklistRequest(deviceId)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.ClearTasklistResponse(reply) : null
    })
}

export function SetPcControl (deviceId, enable) {
  const request = encode.SetPcControlRequest(deviceId, enable)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetPcControlResponse(reply) : null
    })
}

export function SetInterlock (deviceId, interlock) {
  const request = encode.SetInterlockRequest(deviceId, interlock)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetInterlockResponse(reply) : null
    })
}

export function ActivateKeypads (deviceId, reader1, reader2, reader3, reader4) {
  const request = encode.ActivateKeypadsRequest(deviceId, reader1, reader2, reader3, reader4)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.ActivateKeypadsResponse(reply) : null
    })
}

export function SetDoorPasscodes (deviceId, door, passcode1, passcode2, passcode3, passcode4) {
  const request = encode.SetDoorPasscodesRequest(deviceId, door, passcode1, passcode2, passcode3, passcode4)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetDoorPasscodesResponse(reply) : null
    })
}

export function RestoreDefaultParameters (deviceId) {
  const request = encode.RestoreDefaultParametersRequest(deviceId)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.RestoreDefaultParametersResponse(reply) : null
    })
}
