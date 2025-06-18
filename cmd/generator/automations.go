package main

import (
	yaml "gopkg.in/yaml.v3"

	"github.com/angelokurtis/go-home-automations/internal/errors"
)

type Automations []Automation

func UnmarshalAutomations(data []byte) (Automations, error) {
	var r Automations
	if err := yaml.Unmarshal(data, &r); err != nil {
		return nil, errors.WithStack(err)
	}

	return r, nil
}

func (r *Automations) Marshal() ([]byte, error) {
	b, err := yaml.Marshal(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return b, nil
}

type Automation struct {
	ID          string    `yaml:"id,omitempty"`
	Alias       string    `yaml:"alias,omitempty"`
	Description string    `yaml:"description,omitempty"`
	Triggers    []Trigger `yaml:"triggers,omitempty"`
	Actions     []Action  `yaml:"actions,omitempty"`
	Mode        string    `yaml:"mode,omitempty"`
}

type Action struct {
	Choose []Choose `yaml:"choose,omitempty"`
}

type Choose struct {
	Conditions []Condition `yaml:"conditions,omitempty"`
	Sequence   []Sequence  `yaml:"sequence,omitempty"`
}

type Condition struct {
	Condition string `yaml:"condition,omitempty"`
	EntityID  string `yaml:"entity_id,omitempty"`
	State     string `yaml:"state,omitempty"`
}

type Sequence struct {
	Target Target `yaml:"target,omitempty"`
	Action string `yaml:"action,omitempty"`
	Data   any    `yaml:"data"`
}

type Target struct {
	EntityID string `yaml:"entity_id,omitempty"`
}

type Trigger struct {
	EntityID string `yaml:"entity_id,omitempty"`
	Trigger  string `yaml:"trigger,omitempty"`
}
