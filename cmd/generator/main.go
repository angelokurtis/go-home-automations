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
			_, others := lo.Difference([]string{sw}, group.Switches)
			var chooses []Choose
			for _, other := range others {
				chooses = append(chooses, Choose{
					Conditions: []Condition{
						{Condition: "state", EntityID: sw, State: "on"},
						{Condition: "state", EntityID: other, State: "off"},
					},
					Sequence: []Sequence{{
						Target:  Target{EntityID: other},
						Service: "switch.turn_on",
					}},
				})
				chooses = append(chooses, Choose{
					Conditions: []Condition{
						{Condition: "state", EntityID: sw, State: "off"},
						{Condition: "state", EntityID: other, State: "on"},
					},
					Sequence: []Sequence{{
						Target:  Target{EntityID: other},
						Service: "switch.turn_off",
					}},
				})
			}
			automations = append(automations, Automation{
				Alias: fmt.Sprintf("Synchronize %q", sw),
				Triggers: []Trigger{{
					Platform: "state",
					EntityID: sw,
				}},
				Actions: []Action{{Choose: chooses}},
				Mode:    "single",
			})
		}
	}

	y, err := automations.Marshal()
	if err != nil {
		return err
	}

	// Write YAML to file
	outputFile := "automations.yaml"
	if err := os.WriteFile(outputFile, y, 0o644); err != nil {
		return fmt.Errorf("failed to write YAML to file: %w", err)
	}

	slog.InfoContext(ctx, "YAML successfully written to file", "file", outputFile)

	return nil
}
