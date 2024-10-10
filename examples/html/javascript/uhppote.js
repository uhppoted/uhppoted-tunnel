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

export function GetController (controller) {
  const request = encode.GetControllerRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetControllerResponse(reply) : null
    })
}

export function SetIP (controller, address, netmask, gateway) {
  const request = encode.SetIPRequest(controller, address, netmask, gateway)

  return udp.send(request, true)
    .then(() => {
      return true
    })
}

export function GetTime (controller) {
  const request = encode.GetTimeRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetTimeResponse(reply) : null
    })
}

export function SetTime (controller, time) {
  const request = encode.SetTimeRequest(controller, time)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetTimeResponse(reply) : null
    })
}

export function GetStatus (controller) {
  const request = encode.GetStatusRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetStatusResponse(reply) : null
    })
}

export function GetListener (controller) {
  const request = encode.GetListenerRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetListenerResponse(reply) : null
    })
}

export function SetListener (controller, address, port, interval) {
  const request = encode.SetListenerRequest(controller, address, port, interval)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetListenerResponse(reply) : null
    })
}

export function GetDoorControl (controller, door) {
  const request = encode.GetDoorControlRequest(controller, door)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetDoorControlResponse(reply) : null
    })
}

export function SetDoorControl (controller, door, mode, delay) {
  const request = encode.SetDoorControlRequest(controller, door, mode, delay)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetDoorControlResponse(reply) : null
    })
}

export function OpenDoor (controller, door) {
  const request = encode.OpenDoorRequest(controller, door)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.OpenDoorResponse(reply) : null
    })
}

export function GetCards (controller) {
  const request = encode.GetCardsRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetCardsResponse(reply) : null
    })
}

export function GetCard (controller, cardNumber) {
  const request = encode.GetCardRequest(controller, cardNumber)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetCardResponse(reply) : null
    })
}

export function GetCardByIndex (controller, cardIndex) {
  const request = encode.GetCardByIndexRequest(controller, cardIndex)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetCardByIndexResponse(reply) : null
    })
}

export function PutCard (controller, cardNumber, startDate, endDate, door1, door2, door3, door4, PIN) {
  const request = encode.PutCardRequest(controller, cardNumber, startDate, endDate, door1, door2, door3, door4, PIN)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.PutCardResponse(reply) : null
    })
}

export function DeleteCard (controller, cardNumber) {
  const request = encode.DeleteCardRequest(controller, cardNumber)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.DeleteCardResponse(reply) : null
    })
}

export function DeleteAllCards (controller) {
  const request = encode.DeleteCardsRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.DeleteAllCardsResponse(reply) : null
    })
}

export function GetEvent (controller, eventIndex) {
  const request = encode.GetEventRequest(controller, eventIndex)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetEventResponse(reply) : null
    })
}

export function GetEventIndex (controller) {
  const request = encode.GetEventIndexRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetEventIndexResponse(reply) : null
    })
}

export function SetEventIndex (controller, eventIndex) {
  const request = encode.SetEventIndexRequest(controller, eventIndex)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetEventIndexResponse(reply) : null
    })
}

export function RecordSpecialEvents (controller, enable) {
  const request = encode.RecordSpecialEventsRequest(controller, enable)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.RecordSpecialEventsResponse(reply) : null
    })
}

export function GetTimeProfile (controller, profileId) {
  const request = encode.GetTimeProfileRequest(controller, profileId)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetTimeProfileResponse(reply) : null
    })
}

export function SetTimeProfile (controller, profileId, startDate, endDate, monday, tuesday, wednesday, thursday, friday, saturday, sunday, segment1Start, segment1End, segment2Start, segment2End, segment3Start, segment3End, linkedProfileId) {
  const request = encode.SetTimeProfileRequest(controller, profileId, startDate, endDate, monday, tuesday, wednesday, thursday, friday, saturday, sunday, segment1Start, segment1End, segment2Start, segment2End, segment3Start, segment3End, linkedProfileId)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetTimeProfileResponse(reply) : null
    })
}

export function DeleteAllTimeProfiles (controller) {
  const request = encode.DeleteAllTimeProfilesRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.DeleteAllTimeProfilesResponse(reply) : null
    })
}

export function AddTask (controller, startDate, endDate, monday, tuesday, wednesday, thursday, friday, saturday, sunday, startTime, door, taskType, moreCards) {
  const request = encode.AddTaskRequest(controller, startDate, endDate, monday, tuesday, wednesday, thursday, friday, saturday, sunday, startTime, door, taskType, moreCards)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.AddTaskResponse(reply) : null
    })
}

export function RefreshTasklist (controller) {
  const request = encode.RefreshTasklistRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.RefreshTasklistResponse(reply) : null
    })
}

export function ClearTasklist (controller) {
  const request = encode.ClearTasklistRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.ClearTasklistResponse(reply) : null
    })
}

export function SetPcControl (controller, enable) {
  const request = encode.SetPcControlRequest(controller, enable)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetPcControlResponse(reply) : null
    })
}

export function SetInterlock (controller, interlock) {
  const request = encode.SetInterlockRequest(controller, interlock)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetInterlockResponse(reply) : null
    })
}

export function ActivateKeypads (controller, reader1, reader2, reader3, reader4) {
  const request = encode.ActivateKeypadsRequest(controller, reader1, reader2, reader3, reader4)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.ActivateKeypadsResponse(reply) : null
    })
}

export function SetDoorPasscodes (controller, door, passcode1, passcode2, passcode3, passcode4) {
  const request = encode.SetDoorPasscodesRequest(controller, door, passcode1, passcode2, passcode3, passcode4)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetDoorPasscodesResponse(reply) : null
    })
}

export function RestoreDefaultParameters (controller) {
  const request = encode.RestoreDefaultParametersRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.RestoreDefaultParametersResponse(reply) : null
    })
}
