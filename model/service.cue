package model

#Set: {
	[string]: true
}

#Driver: "swarm" | "systemd"
#Type:   "service" | "oneshot"

#Service: {
	driver!:      #Driver
	type!:        #Type
	capabilities: #Set
	options?: _
}
