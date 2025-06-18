package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/samber/lo"
)

func main() {
	ctx := context.Background() // Create base context

	if err := run(ctx); err != nil {
		msg := fmt.Sprintf("Application exited with error: %+v", err)
		slog.ErrorContext(ctx, msg)
		os.Exit(1)
	}

	slog.InfoContext(ctx, "Application exited")
}

func run(ctx context.Context) error {
	automations := make(Automations, 0)

	for _, group := range groups {
		for _, sw := range group.Switches {
			others, _ := lo.Difference([]string{sw}, group.Switches)
			automations = append(automations, Automation{
				Alias: "🔗 " + sw,
				Triggers: []Trigger{{
					EntityID: sw,
					Trigger:  "state",
				}},
				Actions: []Action{{
					Choose: []Choose{
						{
							Conditions: []Condition{{
								Condition: "state",
								EntityID:  sw,
								State:     "on",
							}},
							Sequence: lo.Map(others, func(item string, index int) Sequence {
								return Sequence{
									Target: Target{EntityID: item},
									Action: "switch.turn_on",
								}
							}),
						},
						{
							Conditions: []Condition{{
								Condition: "state",
								EntityID:  sw,
								State:     "off",
							}},
							Sequence: lo.Map(others, func(item string, index int) Sequence {
								return Sequence{
									Target: Target{EntityID: item},
									Action: "switch.turn_off",
								}
							}),
						},
					},
				}},
				Mode: "single",
			})
		}
	}

	return nil
}
