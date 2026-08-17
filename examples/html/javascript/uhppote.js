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

export function GetController(controller) { 
  const request = encode.GetControllerRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetControllerResponse(reply) : null        
    })
}

export function SetIPv4(controller, address, netmask, gateway) { 
  const request = encode.SetIPv4Request(controller, address, netmask, gateway)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetIPv4Response(reply) : null        
    })
}

export function GetTime(controller) { 
  const request = encode.GetTimeRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetTimeResponse(reply) : null        
    })
}

export function SetTime(controller, datetime) { 
  const request = encode.SetTimeRequest(controller, datetime)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetTimeResponse(reply) : null        
    })
}

export function GetListener(controller) { 
  const request = encode.GetListenerRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetListenerResponse(reply) : null        
    })
}

export function SetListener(controller, address, port, interval) { 
  const request = encode.SetListenerRequest(controller, address, port, interval)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetListenerResponse(reply) : null        
    })
}

export function GetDoor(controller, door) { 
  const request = encode.GetDoorRequest(controller, door)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetDoorResponse(reply) : null        
    })
}

export function SetDoor(controller, door, mode, delay) { 
  const request = encode.SetDoorRequest(controller, door, mode, delay)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetDoorResponse(reply) : null        
    })
}

export function OpenDoor(controller, door) { 
  const request = encode.OpenDoorRequest(controller, door)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.OpenDoorResponse(reply) : null        
    })
}

export function GetStatus(controller) { 
  const request = encode.GetStatusRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetStatusResponse(reply) : null        
    })
}

export function GetCards(controller) { 
  const request = encode.GetCardsRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetCardsResponse(reply) : null        
    })
}

export function GetCard(controller, cardNumber) { 
  const request = encode.GetCardRequest(controller, cardNumber)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetCardResponse(reply) : null        
    })
}

export function GetCardAtIndex(controller, index) { 
  const request = encode.GetCardAtIndexRequest(controller, index)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetCardAtIndexResponse(reply) : null        
    })
}

export function PutCard(controller, card, startDate, endDate, door1, door2, door3, door4, PIN, firstCard) { 
  const request = encode.PutCardRequest(controller, card, startDate, endDate, door1, door2, door3, door4, PIN, firstCard)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.PutCardResponse(reply) : null        
    })
}

export function DeleteCard(controller, cardNumber) { 
  const request = encode.DeleteCardRequest(controller, cardNumber)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.DeleteCardResponse(reply) : null        
    })
}

export function DeleteAllCards(controller) { 
  const request = encode.DeleteCardsRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.DeleteAllCardsResponse(reply) : null        
    })
}

export function GetEvent(controller, eventIndex) { 
  const request = encode.GetEventRequest(controller, eventIndex)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetEventResponse(reply) : null        
    })
}

export function GetEventIndex(controller) { 
  const request = encode.GetEventIndexRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetEventIndexResponse(reply) : null        
    })
}

export function SetEventIndex(controller, eventIndex) { 
  const request = encode.SetEventIndexRequest(controller, eventIndex)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetEventIndexResponse(reply) : null        
    })
}

export function RecordSpecialEvents(controller, enabled) { 
  const request = encode.RecordSpecialEventsRequest(controller, enabled)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.RecordSpecialEventsResponse(reply) : null        
    })
}

export function GetTimeProfile(controller, profile) { 
  const request = encode.GetTimeProfileRequest(controller, profile)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetTimeProfileResponse(reply) : null        
    })
}

export function SetTimeProfile(controller, profile, startDate, endDate, monday, tuesday, wednesday, thursday, friday, saturday, sunday, segment1Start, segment1End, segment2Start, segment2End, segment3Start, segment3End, linkedProfileId) { 
  const request = encode.SetTimeProfileRequest(controller, profile, startDate, endDate, monday, tuesday, wednesday, thursday, friday, saturday, sunday, segment1Start, segment1End, segment2Start, segment2End, segment3Start, segment3End, linkedProfileId)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetTimeProfileResponse(reply) : null        
    })
}

export function ClearTimeProfiles(controller) { 
  const request = encode.ClearTimeProfilesRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.ClearTimeProfilesResponse(reply) : null        
    })
}

export function AddTask(controller, task, startDate, endDate, monday, tuesday, wednesday, thursday, friday, saturday, sunday, startTime, door, moreCards) { 
  const request = encode.AddTaskRequest(controller, task, startDate, endDate, monday, tuesday, wednesday, thursday, friday, saturday, sunday, startTime, door, moreCards)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.AddTaskResponse(reply) : null        
    })
}

export function RefreshTaskList(controller) { 
  const request = encode.RefreshTaskListRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.RefreshTaskListResponse(reply) : null        
    })
}

export function ClearTaskList(controller) { 
  const request = encode.ClearTaskListRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.ClearTaskListResponse(reply) : null        
    })
}

export function SetPCControl(controller, enabled) { 
  const request = encode.SetPCControlRequest(controller, enabled)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetPCControlResponse(reply) : null        
    })
}

export function SetInterlock(controller, interlock) { 
  const request = encode.SetInterlockRequest(controller, interlock)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetInterlockResponse(reply) : null        
    })
}

export function ActivateKeypads(controller, reader1, reader2, reader3, reader4) { 
  const request = encode.ActivateKeypadsRequest(controller, reader1, reader2, reader3, reader4)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.ActivateKeypadsResponse(reply) : null        
    })
}

export function SetDoorPasscodes(controller, door, passcode1, passcode2, passcode3, passcode4) { 
  const request = encode.SetDoorPasscodesRequest(controller, door, passcode1, passcode2, passcode3, passcode4)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetDoorPasscodesResponse(reply) : null        
    })
}

export function GetAntiPassback(controller) { 
  const request = encode.GetAntipassbackRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.GetAntiPassbackResponse(reply) : null        
    })
}

export function SetAntiPassback(controller, antipassback) { 
  const request = encode.SetAntipassbackRequest(controller, antipassback)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetAntiPassbackResponse(reply) : null        
    })
}

export function SetFirstCard(controller, door, startTime, endTime, activeMode, inactiveMode, monday, tuesday, wednesday, thursday, friday, saturday, sunday) { 
  const request = encode.SetFirstCardRequest(controller, door, startTime, endTime, activeMode, inactiveMode, monday, tuesday, wednesday, thursday, friday, saturday, sunday)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.SetFirstCardResponse(reply) : null        
    })
}

export function RestoreDefaultParameters(controller) { 
  const request = encode.RestoreDefaultParametersRequest(controller)

  return udp.send(request)
    .then(reply => {
      return reply ? decode.RestoreDefaultParametersResponse(reply) : null        
    })
}



