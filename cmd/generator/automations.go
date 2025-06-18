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
	ID          string    `yaml:"id"`
	Alias       string    `yaml:"alias"`
	Description string    `yaml:"description"`
	Triggers    []Trigger `yaml:"triggers"`
	Actions     []Action  `yaml:"actions"`
	Mode        string    `yaml:"mode"`
}

type Action struct {
	Choose []Choose `yaml:"choose"`
}

type Choose struct {
	Conditions []Condition `yaml:"conditions"`
	Sequence   []Sequence  `yaml:"sequence"`
}

type Condition struct {
	Condition string `yaml:"condition"`
	EntityID  string `yaml:"entity_id"`
	State     string `yaml:"state"`
}

type Sequence struct {
	Target Target `yaml:"target"`
	Action string `yaml:"action"`
}

type Target struct {
	EntityID string `yaml:"entity_id"`
}

type Trigger struct {
	EntityID string `yaml:"entity_id"`
	Trigger  string `yaml:"trigger"`
}
