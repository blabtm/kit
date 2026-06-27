package svc

// Driver represents execution environment type.
// 
// `docker` - service will run in container.
// `system` - service will run natively.
// `firmware` - service will run on embedded device.
#Driver: "docker" | "system" | "firmware"

// Capability is an optional service functionality.
//
// `notify` - service can be notified to reload it's state.
//
//    Notification type depends on service driver:
//      * `docker` - signal will be sent to service's container.
//      * `system` - signal will be sent to service's process.
//      * `firmware` - new state will be sent to service's control channel.
//    
//    If the service does not support this capability, it'll be reloaded
//    with the new state (convinient for highly-constrained devices with
//    embedded configuration).
#Capability: "notify"

#Version: string & =~"^[0-9]+[.][0-9]+[.][0-9]+$"

#Service: {
  // version indicates service's version.
  version!: #Version

  // driver indicates service's execution environment type.
	driver!: #Driver

  // capabilities specifies the set of optional functionalities.
	capabilities: {
    [#Capability]: true
  }
}
