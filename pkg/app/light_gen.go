// Code generated by ifacemaker; DO NOT EDIT.

package app

// Light defines the interface for controlling a light entity
type Light interface {
	// TurnOn a light entity. Takes an entityId and an optional
	// map that is translated into service_data.
	TurnOn(entityId string, serviceData ...map[string]any)
	// Toggle a light entity. Takes an entityId and an optional
	// map that is translated into service_data.
	Toggle(entityId string, serviceData ...map[string]any)
	TurnOff(entityId string)
}
