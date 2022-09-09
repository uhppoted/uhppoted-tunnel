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

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.GetControllerResponse(replies[0])
      }

      return null
    })
}

export function SetIP (deviceId, address, netmask, gateway) {
  const request = encode.SetIPRequest(deviceId, address, netmask, gateway)

  return udp.send(request, '0.1ms')
    .then(replies => {
      return true
    })
}

export function GetTime (deviceId) {
  const request = encode.GetTimeRequest(deviceId)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.GetTimeResponse(replies[0])
      }

      return null
    })
}

export function SetTime (deviceId, time) {
  const request = encode.SetTimeRequest(deviceId, time)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.SetTimeResponse(replies[0])
      }

      return null
    })
}

export function GetStatus (deviceId) {
  const request = encode.GetStatusRequest(deviceId)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.GetStatusResponse(replies[0])
      }

      return null
    })
}

export function GetListener (deviceId) {
  const request = encode.GetListenerRequest(deviceId)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.GetListenerResponse(replies[0])
      }

      return null
    })
}

export function SetListener (deviceId, address, port) {
  const request = encode.SetListenerRequest(deviceId, address, port)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.SetListenerResponse(replies[0])
      }

      return null
    })
}

export function GetDoorControl (deviceId, door) {
  const request = encode.GetDoorControlRequest(deviceId, door)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.GetDoorControlResponse(replies[0])
      }

      return null
    })
}

export function SetDoorControl (deviceId, door, mode, delay) {
  const request = encode.SetDoorControlRequest(deviceId, door, mode, delay)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.SetDoorControlResponse(replies[0])
      }

      return null
    })
}

export function OpenDoor (deviceId, door) {
  const request = encode.OpenDoorRequest(deviceId, door)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.OpenDoorResponse(replies[0])
      }

      return null
    })
}

export function GetCards (deviceId) {
  const request = encode.GetCardsRequest(deviceId)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.GetCardsResponse(replies[0])
      }

      return null
    })
}

export function GetCard (deviceId, cardNumber) {
  const request = encode.GetCardRequest(deviceId, cardNumber)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.GetCardResponse(replies[0])
      }

      return null
    })
}

export function GetCardByIndex (deviceId, cardIndex) {
  const request = encode.GetCardByIndexRequest(deviceId, cardIndex)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.GetCardByIndexResponse(replies[0])
      }

      return null
    })
}

export function PutCard (deviceId, cardNumber, startDate, endDate, door1, door2, door3, door4) {
  const request = encode.PutCardRequest(deviceId, cardNumber, startDate, endDate, door1, door2, door3, door4)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.PutCardResponse(replies[0])
      }

      return null
    })
}

export function DeleteCard (deviceId, cardNumber) {
  const request = encode.DeleteCardRequest(deviceId, cardNumber)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.DeleteCardResponse(replies[0])
      }

      return null
    })
}

export function DeleteAllCards (deviceId) {
  const request = encode.DeleteCardsRequest(deviceId)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.DeleteAllCardsResponse(replies[0])
      }

      return null
    })
}

export function GetEvent (deviceId, eventIndex) {
  const request = encode.GetEventRequest(deviceId, eventIndex)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.GetEventResponse(replies[0])
      }

      return null
    })
}

export function GetEventIndex (deviceId) {
  const request = encode.GetEventIndexRequest(deviceId)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.GetEventIndexResponse(replies[0])
      }

      return null
    })
}

export function SetEventIndex (deviceId, eventIndex) {
  const request = encode.SetEventIndexRequest(deviceId, eventIndex)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.SetEventIndexResponse(replies[0])
      }

      return null
    })
}

export function RecordSpecialEvents (deviceId, enable) {
  const request = encode.RecordSpecialEventsRequest(deviceId, enable)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.RecordSpecialEventsResponse(replies[0])
      }

      return null
    })
}

export function GetTimeProfile (deviceId, profileId) {
  const request = encode.GetTimeProfileRequest(deviceId, profileId)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.GetTimeProfileResponse(replies[0])
      }

      return null
    })
}

export function SetTimeProfile (deviceId, profileId, startDate, endDate, monday, tuesday, wednesday, thursday, friday, saturday, sunday, segment1Start, segment1End, segment2Start, segment2End, segment3Start, segment3End, linkedProfileId) {
  const request = encode.SetTimeProfileRequest(deviceId, profileId, startDate, endDate, monday, tuesday, wednesday, thursday, friday, saturday, sunday, segment1Start, segment1End, segment2Start, segment2End, segment3Start, segment3End, linkedProfileId)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.SetTimeProfileResponse(replies[0])
      }

      return null
    })
}

export function DeleteAllTimeProfiles (deviceId) {
  const request = encode.DeleteAllTimeProfilesRequest(deviceId)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.DeleteAllTimeProfilesResponse(replies[0])
      }

      return null
    })
}

export function AddTask (deviceId, startDate, endDate, monday, tuesday, wednesday, thursday, friday, saturday, sunday, startTime, door, taskType, moreCards) {
  const request = encode.AddTaskRequest(deviceId, startDate, endDate, monday, tuesday, wednesday, thursday, friday, saturday, sunday, startTime, door, taskType, moreCards)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.AddTaskResponse(replies[0])
      }

      return null
    })
}

export function RefreshTasklist (deviceId) {
  const request = encode.RefreshTasklistRequest(deviceId)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.RefreshTasklistResponse(replies[0])
      }

      return null
    })
}

export function ClearTasklist (deviceId) {
  const request = encode.ClearTasklistRequest(deviceId)

  return udp.send(request, '0s')
    .then(replies => {
      if (replies.length > 0) {
        return decode.ClearTasklistResponse(replies[0])
      }

      return null
    })
}
