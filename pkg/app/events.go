package app

import (
	"encoding/json"
)

func (r *StateChangeEvent) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type StateChangeEvent struct {
	Type     string
	EntityID string
	OldState string
	NewState string
}
