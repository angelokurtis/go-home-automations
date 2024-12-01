package app

import "time"

import "encoding/json"

type Events []Event

func UnmarshalEvents(data []byte) (Events, error) {
	var r Events
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Events) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Event struct {
	Event EventClass `json:"event"`
	ID    int64      `json:"id"`
	Type  string     `json:"type"`
}

type EventClass struct {
	Context   Context   `json:"context"`
	Data      Data      `json:"data"`
	EventType string    `json:"event_type"`
	Origin    string    `json:"origin"`
	TimeFired time.Time `json:"time_fired"`
}

type Context struct {
	ID       string  `json:"id"`
	ParentID *string `json:"parent_id"`
	UserID   *string `json:"user_id"`
}

type Data struct {
	EntityID string      `json:"entity_id"`
	NewState EntityState `json:"new_state"`
	OldState EntityState `json:"old_state"`
}

type EntityState struct {
	Attributes   map[string]string `json:"attributes"`
	Context      Context           `json:"context"`
	EntityID     string            `json:"entity_id"`
	LastChanged  time.Time         `json:"last_changed"`
	LastReported time.Time         `json:"last_reported"`
	LastUpdated  time.Time         `json:"last_updated"`
	EntityState  string            `json:"state"`
}
