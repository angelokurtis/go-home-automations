package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
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
			automations = append(automations, Automation{
				Alias:       "Sync " + sw,
				Description: "",
				Triggers:    nil,
				Actions:     nil,
				Mode:        "",
			})
		}
	}
	return nil
}
