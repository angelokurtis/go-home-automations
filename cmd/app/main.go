package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/angelokurtis/go-otel/starter"
	"github.com/lmittmann/tint"
	"go.uber.org/automaxprocs/maxprocs"
	"golang.org/x/sync/errgroup"

	"github.com/angelokurtis/go-home-automations/internal/errors"
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

// run manages app lifecycle, signal handling, and runner.
func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		AddSource:  true,
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
	})))

	_, shutdown, err := starter.StartProviders(context.Background())
	if err != nil {
		return errors.WithStack(err)
	}
	defer shutdown()

	// Set up GOMAXPROCS to utilize available CPU cores
	undo, err := maxprocs.Set()
	defer undo()

	if err != nil {
		return errors.Errorf("failed to set GOMAXPROCS: %w", err)
	}

	runner, cleanup, err := NewRunner(ctx)
	if err != nil {
		return err
	}
	defer cleanup()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		defer stop()
		return runner.Run(ctx)
	})

	g.Go(func() error {
		<-ctx.Done()
		return nil
	})

	return g.Wait()
}
