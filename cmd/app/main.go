package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lmittmann/tint"

	"github.com/angelokurtis/go-home-automations/internal/maxprocs"
)

func main() {
	// Create a context that can be canceled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancel is called to release resources

	// Set up logging
	logger := slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		AddSource:  true,
		Level:      slog.LevelInfo,
		TimeFormat: time.Kitchen,
	}))
	slog.SetDefault(logger)

	// Set up GOMAXPROCS to utilize available CPU cores
	_, undo := maxprocs.SetUp(logger)
	defer undo()

	// Initialize the application
	appRunner, cleanup, err := newAppRunner(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to initialize application", tint.Err(err))
		os.Exit(1)
	}

	// Channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// Channel to report errors from the process
	errChan := make(chan error, 1)

	// Start the main application logic in a new goroutine
	go func() { errChan <- appRunner.Run(ctx) }()

	select {
	case err = <-errChan:
		// App completed with error
		if err != nil {
			slog.ErrorContext(ctx, "Application finished with error", tint.Err(err))
			cleanup()
			os.Exit(1)
		}

		slog.DebugContext(ctx, "Application completed successfully")
	case sig := <-sigChan:
		// Received a termination signal
		slog.WarnContext(ctx, "Received termination signal", slog.String("signal", sig.String()))
		cancel()
	}

	// Run cleanup
	slog.DebugContext(ctx, "Running cleanup")
	cleanup()
}
